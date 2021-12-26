# Cole
Cole can use his sixth sense to tell secrets about your dashboards usage.

# Contributing
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