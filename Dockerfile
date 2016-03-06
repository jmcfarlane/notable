FROM debian
MAINTAINER John McFarlane

EXPOSE 8082
RUN apt-get update && apt-get install -y python-pip
RUN pip install notable==0.3.1
ENTRYPOINT ["notable", "-f", "-b 0.0.0.0"]
