#!/bin/bash

set -a
source "xlog.env"

make all #test

#export LOG_GOID=1
#export LOG_SRC=1
#export LOG_SRC_PKG=1
export LOG_SRC_EXT=1
#export LOG_ID=1
#export LOG_SUM_FULL=1
#export LOG_TIME=1
#export LOG_SUM_ALONE=0
export LOG_RATE_LIMIT=1
#export LOG_FORMAT=std
#export LOG_FORMAT=text
export LOG_FORMAT=tint
#export LOG_FORMAT=json

mkdir -p logs

./xlogscan -log-format json test | tee logs/test.log | ./xlogscan

#LOG_SUM_CHAIN=1 ./xlogscan -log-format prod test | tee logs/test.log | ./xlogscan --chain
