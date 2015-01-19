# oplog-dump worker
FROM ubuntu:14.04
RUN apt-get update
RUN apt-get install -y wget build-essential

# Golang
RUN apt-get install -y git golang bzr mercurial bash
RUN GOPATH=/etc/go go get launchpad.net/godeb
RUN apt-get remove -y golang golang-go golang-doc golang-src
RUN /etc/go/bin/godeb install 1.2.1

# Mongo
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
RUN echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | sudo tee /etc/apt/sources.list.d/mongodb.list
RUN apt-get -y update
RUN apt-get install -y mongodb-org

# Oplog dump
RUN mkdir -p /etc/go/src /github.com/Clever/oplog-dump
ADD . /etc/go/src/github.com/Clever/oplog-dump
RUN GOPATH=/etc/go go get github.com/Clever/oplog-dump/...
RUN GOPATH=/etc/go go build -o /usr/local/bin/oplogdump github.com/Clever/oplog-dump/cmd/oplog-dump

# Gearcmd
RUN curl -L https://github.com/Clever/gearcmd/releases/download/v0.3.3/gearcmd-v0.3.3-linux-amd64.tar.gz | tar xz -C /usr/local/bin --strip-components 1

CMD ["/etc/go/src/github.com/Clever/oplog-dump/run_as_worker.sh"]
