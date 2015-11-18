FROM google/debian:wheezy

RUN apt-get -y update && apt-get -y install curl && curl -L https://github.com/Clever/gearcmd/releases/download/v0.4.0/gearcmd-v0.4.0-linux-amd64.tar.gz | tar xz -C /usr/local/bin --strip-components 1

# Install Mongo
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10		
RUN echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | sudo tee /etc/apt/sources.list.d/mongodb.list		
RUN apt-get -y update && apt-get install -y mongodb-org

COPY bin/oplog-dump /usr/bin/oplog-dump

CMD ["gearcmd", "--name", "oplog-dump", "--cmd", "/usr/bin/oplog-dump", "--cmdtimeout", "1h"]
