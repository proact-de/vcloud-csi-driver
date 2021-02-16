# vcloud-csi-driver

[![Go Reference](https://pkg.go.dev/badge/github.com/proact-de/vcloud-csi-driver.svg)](https://pkg.go.dev/github.com/proact-de/vcloud-csi-driver) [![Go Report Card](https://goreportcard.com/badge/github.com/proact-de/vcloud-csi-driver)](https://goreportcard.com/report/github.com/proact-de/vcloud-csi-driver) [![Docker Size](https://img.shields.io/docker/image-size/proactcloud/vcloud-csi-driver/latest)](https://hub.docker.com/r/proactcloud/vcloud-csi-driver) [![Docker Pulls](https://img.shields.io/docker/pulls/proactcloud/vcloud-csi-driver)](https://hub.docker.com/r/proactcloud/vcloud-csi-driver)

A Container Storage Interface ([CSI](https://github.com/container-storage-interface/spec)) Driver for [VMWare vCloud Director](https://www.vmware.com/de/products/cloud-director.html). The CSI plugin allows you to use independant disks with your preferred Container Orchestrator. If you need guidance how to install this take a look at our [documentation](https://proact-de.github.io/vcloud-csi-driver/).

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html).

```console
git clone https://github.com/proact-de/vcloud-csi-driver.git
cd vcloud-csi-driver

make generate build

./bin/vcloud-csi-driver -h
```

## Security

If you find a security issue please contact devops@proactcloud.de first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

* [Thomas Boerger](https://github.com/tboerger)

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2021 Proact Deutschland GmbH <devops@proactcloud.de>
```
