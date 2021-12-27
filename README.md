[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/cole)](https://artifacthub.io/packages/search?repo=cole)
[![Build](https://github.com/nicolastakashi/cole/actions/workflows/docker-publish.yml/badge.svg?event=branch_protection_rule)](https://github.com/nicolastakashi/cole/actions/workflows/docker-publish.yml)

# Cole
Cole can use his sixth sense to give you metrics about your Grafana dashboards

## Overview

Cole is a lightweight service that handles HTTP logs of Grafana to provide insights about the usage of Grafana instances and expose it through Prometheus metrics.

## Grafana router logging dependency
Cole uses its sixth sense by Grafana HTTP logs, and because of this, you need to enable the `router_logging` to log all HTTP requests (not just errors).

For more information about router logging, please, check the [Grafana official documentation](https://grafana.com/docs/grafana/latest/administration/configuration/#router_logging).

## Metrics Documentation

See the [docs](./docs/README.md) directory for more information on the exposed metrics.

## Contributing
Contributions are very welcome! See our [CONTRIBUTING.md](CONTRIBUTING.md) for more information.

## Docker images

Docker images are available on [Docker Hub](https://hub.docker.com/r/ntakashi/cole).

## Building from source

To build Cole from source code, first ensure that you have a working
Go environment with [version 1.16 or greater installed](https://golang.org/doc/install).

To build the source code you can use the `make build`, which will compile in
the assets so that Cole can be run from anywhere:

```bash
$ mkdir -p $GOPATH/src/github.com/nicolastakashi/cole
$ cd $GOPATH/src/github.com/nicolastakashi/cole
$ git clone https://github.com/nicolastakashi/cole.git
$ cd cole
$ make build
$ ./cole server <args>
```

The Makefile provides several targets:

  * *build*: build the `cole`
  * *fmt*: format the source code
  * *vet*: check the source code for common errors
  * *tests*: run unit tests