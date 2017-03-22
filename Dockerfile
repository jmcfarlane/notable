FROM alpine
MAINTAINER John McFarlane

EXPOSE 8080

# Add the compiled binary
COPY notable /notable

CMD ["/notable", "-daemon=false", "-browser=false", "-bind=0.0.0.0"]
