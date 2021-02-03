package metrics

import (
	"context"
	"io"
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/proact-de/vcloud-csi-driver/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// Server simply defines a webserver to handle metrics exporting.
type Server struct {
	registry *prometheus.Registry
	grpc     *grpc_prometheus.ServerMetrics
	http     *http.Server
}

// NewServer initializes a new server to handle metrics exporting.
func NewServer() *Server {
	metrics := &Server{
		registry: prometheus.NewRegistry(),
		grpc:     grpc_prometheus.NewServerMetrics(),
	}

	metrics.grpc.EnableHandlingTimeHistogram()

	metrics.registry.MustRegister(
		metrics.grpc,
	)

	metrics.registry.MustRegister(
		prometheus.NewGoCollector(),
	)

	metrics.registry.MustRegister(
		version.NewCollector(),
	)

	metrics.registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{
		Namespace: "vcloud_csi",
	}))

	return metrics
}

// UnaryServerInterceptor returns an intercaptor for the GRPC server to gather metrics.
func (s *Server) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return s.grpc.UnaryServerInterceptor()
}

// InitializeMetrics initializes all metrics for all gRPC methods registered on a GRPC server.
func (s *Server) InitializeMetrics(server *grpc.Server) {
	s.grpc.InitializeMetrics(server)
}

// ListenAndServe simply wraps the http router for the metrics endpoint.
func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()

	mux.Handle("/", promhttp.HandlerFor(
		s.registry,
		promhttp.HandlerOpts{},
	))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		io.WriteString(w, http.StatusText(http.StatusOK))
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		io.WriteString(w, http.StatusText(http.StatusOK))
	})

	s.http = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return s.http.ListenAndServe()
}

// Shutdown simply wraps the http router for the metrics endpoint.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
