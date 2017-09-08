# oplog-dump worker
FROM ubuntu:14.04

# install gearcmd
RUN apt-get -y update && \
    apt-get install -y -q curl && \
    curl -L https://github.com/Clever/gearcmd/releases/download/0.8.7/gearcmd-v0.8.7-linux-amd64.tar.gz | tar xz -C /usr/local/bin --strip-components 1

# Mongo tooling
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10 && \
    echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | sudo tee /etc/apt/sources.list.d/mongodb.list && \
    apt-get -y update && \
    apt-get install -y mongodb-org

COPY bin/oplog-dump /usr/local/bin/oplog-dump
CMD ["gearcmd", "--name", "oplog-dump", "--cmd", "/usr/local/bin/oplog-dump", "--cmdtimeout", "1h"]
