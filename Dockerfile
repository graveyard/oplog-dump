# oplog-dump worker
FROM ubuntu:14.04

RUN apt-get -y update && \
    apt-get install -y wget build-essential && \
    apt-get -y update && \
    apt-get install -y -q curl

# Mongo
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
RUN echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | sudo tee /etc/apt/sources.list.d/mongodb.list
RUN apt-get -y update && apt-get install -y mongodb-org

COPY bin/oplogdump /usr/local/bin/oplogdump

# Gearcmd
RUN curl -L https://github.com/Clever/gearcmd/releases/download/v0.4.0/gearcmd-v0.4.0-linux-amd64.tar.gz | tar xz -C /usr/local/bin --strip-components 1

CMD ["gearcmd", "--name", "oplog-dump", "--cmd", "/usr/local/bin/oplogdump", "--cmdtimeout", "1h"]
