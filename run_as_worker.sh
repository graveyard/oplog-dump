#!/bin/bash
taskwrapper --name oplog-dump --cmd /usr/local/bin/oplogdump --gearman-host $GEARMAN_HOST --gearman-port $GEARMAN_PORT
