# oplog-dump worker
FROM ubuntu:14.04
RUN apt-get update
RUN apt-get install -y wget build-essential

# Golang
RUN apt-get install -y git golang bzr mercurial bash
RUN GOPATH=/etc/go go get launchpad.net/godeb
RUN apt-get remove -y golang golang-go golang-doc golang-src
RUN /etc/go/bin/godeb install 1.2.1

# Oplog dump
RUN mkdir -p /etc/go/src /github.com/Clever/oplog-dump
ADD . /etc/go/src/github.com/Clever/oplog-dump
RUN GOPATH=/etc/go go get github.com/Clever/oplog-dump/...
RUN GOPATH=/etc/go go build -o /usr/local/bin/oplogdump github.com/Clever/oplog-dump/cmd/oplog-dump

# Gearcmd
RUN wget https://github.com/Clever/gearcmd/releases/download/v0.1.2/gearcmd-v0.1.2-linux-amd64.tar.gz
RUN tar -xvf gearcmd-v0.1.2-linux-amd64.tar.gz
RUN cp gearcmd-v0.1.2-linux-amd64/gearcmd /usr/local/bin/

CMD ["/etc/go/src/github.com/Clever/oplog-dump/run_as_worker.sh"]
