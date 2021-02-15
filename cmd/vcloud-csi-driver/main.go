package main

import (
	"context"
	"errors"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/oklog/run"
	"github.com/proact-de/vcloud-csi-driver/pkg/config"
	"github.com/proact-de/vcloud-csi-driver/pkg/metrics"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/controller"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/identity"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/mount"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/node"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/resize"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/volume"
	"github.com/proact-de/vcloud-csi-driver/pkg/vcloud"
	"github.com/proact-de/vcloud-csi-driver/pkg/version"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	app := &cli.App{
		Name:    "vcloud-csi-driver",
		Version: version.String,
		Usage:   "CSI driver for vCloud Director",
		Authors: []*cli.Author{
			{
				Name:  "Thomas Boerger",
				Email: "thomas.boerger@proact.de",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Value:       "info",
				Usage:       "set logging level",
				EnvVars:     []string{"VCLOUD_CSI_LOG_LEVEL"},
				Destination: &cfg.Logs.Level,
			},
			&cli.BoolFlag{
				Name:        "log-pretty",
				Value:       true,
				Usage:       "enable pretty logging",
				EnvVars:     []string{"VCLOUD_CSI_LOG_PRETTY"},
				Destination: &cfg.Logs.Pretty,
			},
			&cli.BoolFlag{
				Name:        "log-color",
				Value:       true,
				Usage:       "enable colored logging",
				EnvVars:     []string{"VCLOUD_CSI_LOG_COLOR"},
				Destination: &cfg.Logs.Color,
			},
			&cli.StringFlag{
				Name:        "kube-nodename",
				Value:       "",
				Usage:       "Name of the node running on",
				EnvVars:     []string{"VCLOUD_CSI_NODENAME"},
				Destination: &cfg.Kubernetes.Nodename,
			},
			&cli.StringFlag{
				Name:        "kube-namespace",
				Value:       "",
				Usage:       "Namespace running on Kubernetes",
				EnvVars:     []string{"VCLOUD_CSI_NAMESPACE"},
				Destination: &cfg.Kubernetes.Namespace,
			},
			&cli.StringFlag{
				Name:        "kube-podip",
				Value:       "",
				Usage:       "IP address of the running pod",
				EnvVars:     []string{"VCLOUD_CSI_PODIP"},
				Destination: &cfg.Kubernetes.PodIP,
			},
			&cli.StringFlag{
				Name:        "vcloud-href",
				Value:       "",
				Usage:       "URL to access vCloud Director API",
				EnvVars:     []string{"VCLOUD_CSI_HREF"},
				Destination: &cfg.Driver.Href,
			},
			&cli.BoolFlag{
				Name:        "vcloud-insecure",
				Value:       false,
				Usage:       "Skip SSL verify for vCloud Director",
				EnvVars:     []string{"VCLOUD_CSI_INSECURE"},
				Destination: &cfg.Driver.Insecure,
			},
			&cli.StringFlag{
				Name:        "vcloud-username",
				Value:       "",
				Usage:       "Username for vCloud Director",
				EnvVars:     []string{"VCLOUD_CSI_USERNAME"},
				Destination: &cfg.Driver.Username,
			},
			&cli.StringFlag{
				Name:        "vcloud-password",
				Value:       "",
				Usage:       "Password for vCloud Director",
				EnvVars:     []string{"VCLOUD_CSI_PASSWORD"},
				Destination: &cfg.Driver.Password,
			},
			&cli.StringFlag{
				Name:        "vcloud-org",
				Value:       "",
				Usage:       "Organization for vCloud Director",
				EnvVars:     []string{"VCLOUD_CSI_ORG"},
				Destination: &cfg.Driver.Organization,
			},
			&cli.StringFlag{
				Name:        "vcloud-vdc",
				Value:       "",
				Usage:       "VDCs for vCloud Director",
				EnvVars:     []string{"VCLOUD_CSI_VDC"},
				Destination: &cfg.Driver.Datacenter,
			},
			&cli.StringFlag{
				Name:        "csi-endpoint",
				Value:       "unix:///csi/csi.sock",
				Usage:       "Path to unix socket endpoint",
				EnvVars:     []string{"VCLOUD_CSI_ENDOINT"},
				Destination: &cfg.Driver.Endpoint,
			},
		},
		Before: func(c *cli.Context) error {
			setupLogger(cfg)

			cfg.Driver.Endpoint = strings.TrimPrefix(
				cfg.Driver.Endpoint,
				"unix://",
			)

			if !checkEndpointDefined(cfg) {
				os.Exit(1)
			}

			if !ensureSocketRemoved(cfg) {
				os.Exit(1)
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			log.Info().
				Str("version", version.String).
				Str("revision", version.Revision).
				Str("date", version.Date).
				Str("go", version.Go).
				Msg("vCloud Director CSI driver")

			clientHref, err := url.Parse(cfg.Driver.Href)

			if err != nil {
				log.Error().
					Err(err).
					Msg("Failed to parse vCloud Director URL")

				return errors.New("failed to parse vcloud director url")
			}

			client, err := vcloud.NewClient(
				vcloud.WithHref(clientHref),
				vcloud.WithInsecure(cfg.Driver.Insecure),
				vcloud.WithUsername(cfg.Driver.Username),
				vcloud.WithPassword(cfg.Driver.Password),
				vcloud.WithOrganization(cfg.Driver.Organization),
				vcloud.WithDatacenter(cfg.Driver.Datacenter),
			)

			if err != nil {
				log.Error().
					Err(err).
					Msg("Failed to init vCloud Director client")

				return errors.New("failed to init vcloud director client")
			}

			disks, err := client.List()

			if err != nil {
				log.Error().
					Err(err).
					Msg("")
			}

			log.Info().
				Interface("disks", disks).
				Msg("")

			os.Exit(1)

			volumeService := volume.NewService(
				volume.WithClient(client),
			)

			ctrlService := controller.NewService(
				controller.WithServer(cfg.Kubernetes.Nodename),
				controller.WithDatacenter(cfg.Driver.Datacenter),
				controller.WithVolume(volumeService),
			)

			nodeService := node.NewService(
				node.WithServer(cfg.Kubernetes.Nodename),
				node.WithDatacenter(cfg.Driver.Datacenter),
				node.WithVolume(volumeService),
				node.WithMount(mount.NewService()),
				node.WithResize(resize.NewService()),
			)

			identityService := identity.NewService()

			listener, err := net.Listen(
				"unix",
				cfg.Driver.Endpoint,
			)

			if err != nil {
				log.Error().
					Err(err).
					Str("endpoint", cfg.Driver.Endpoint).
					Msg("Failed to create listener")

				os.Exit(1)
			}

			var gr run.Group

			metricsServer := metrics.NewServer()

			grpcServer := grpc.NewServer(
				grpc.UnaryInterceptor(
					grpc_middleware.ChainUnaryServer(
						func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
							log.Debug().
								Interface("req", req).
								Msg("Handling request")

							resp, err := handler(ctx, req)

							if err != nil {
								log.Error().
									Err(err).
									Interface("req", req).
									Interface("resp", resp).
									Msg("Handling failed")

								return resp, err
							}

							log.Debug().
								Msg("Handling finished")

							return resp, err
						},
						metricsServer.UnaryServerInterceptor(),
					),
				),
			)

			metricsServer.InitializeMetrics(grpcServer)

			csi.RegisterControllerServer(grpcServer, ctrlService)
			csi.RegisterNodeServer(grpcServer, nodeService)
			csi.RegisterIdentityServer(grpcServer, identityService)

			gr.Add(func() error {
				log.Info().
					Str("server", "grpc").
					Msg("Starting server")

				return grpcServer.Serve(listener)
			}, func(reason error) {
				grpcServer.GracefulStop()

				log.Info().
					Err(reason).
					Str("server", "grpc").
					Msg("Shutdown gracefully")
			})

			gr.Add(func() error {
				log.Info().
					Str("server", "metrics").
					Msg("Starting server")

				return metricsServer.ListenAndServe()
			}, func(reason error) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := metricsServer.Shutdown(ctx); err != nil {
					log.Info().
						Err(err).
						Str("server", "metrics").
						Msg("Shutdown failed")

					return
				}

				log.Info().
					Err(reason).
					Str("server", "metrics").
					Msg("Shutdown gracefully")
			})

			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(
						stop,
						os.Interrupt,
						syscall.SIGTERM,
						syscall.SIGINT,
					)

					<-stop

					return nil
				}, func(err error) {
					close(stop)
				})
			}

			return gr.Run()
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show the help, so what you see here now",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the current version of that tool",
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
