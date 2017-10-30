FROM golang:1.9.2

# Build time variables (set inside Makefile)
ARG CHOWN_GID=unset
ARG CHOWN_UID=unset
ARG VERSION=unset

ENV GOPATH /go
ENV PATH $PATH:/go/bin
ENV PROJECT /go/src/github.com/jmcfarlane/notable
ENV VERSION=$VERSION

ADD . $PROJECT
WORKDIR $PROJECT

RUN apt-get update && apt-get -y install imagemagick icnsutils zip
RUN make build
