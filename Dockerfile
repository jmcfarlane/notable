FROM scratch
MAINTAINER John McFarlane

EXPOSE 8080

# Add the compiled binary
COPY target/notable-v*.linux-amd64/notable /notable

ENTRYPOINT ["/notable", "-daemon=false", "-browser=false", "-bind=0.0.0.0"]
