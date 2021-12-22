# syntax=docker/dockerfile:1

FROM golang:latest as build

WORKDIR /go/src/github.com/nicolastakashi/cole

RUN apt-get update
RUN useradd -ms /bin/bash cole

COPY --chown=cole:cole . .

RUN make all

RUN pwd
RUN ls -la /go/src/github.com/nicolastakashi/cole

FROM golang:latest as runtime

RUN useradd -ms /bin/bash cole

WORKDIR /cole

COPY --from=build /go/src/github.com/nicolastakashi/cole/bin/* /bin/

USER cole

ENTRYPOINT [ "/bin/cole" ]