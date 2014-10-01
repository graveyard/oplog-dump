#!/bin/bash
gearcmd --name oplog-dump --cmd /usr/local/bin/oplogdump --host $GEARMAN_HOST --port $GEARMAN_PORT
