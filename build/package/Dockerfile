# syntax=docker/dockerfile:1

FROM golang:latest as build

WORKDIR /go/src/github.com/nicolastakashi/cole

RUN apt-get update
RUN useradd -ms /bin/bash cole

COPY --chown=cole:cole . .

RUN make all

FROM gcr.io/distroless/static:latest-amd64

WORKDIR /cole

COPY --from=build /go/src/github.com/nicolastakashi/cole/bin/* /bin/

USER nobody

ENTRYPOINT [ "/bin/cole" ]