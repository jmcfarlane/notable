FROM golang
MAINTAINER John McFarlane

EXPOSE 8080
RUN go get -u github.com/jmcfarlane/notable
ENTRYPOINT ["/go/bin/notable", "-daemon=false"]
