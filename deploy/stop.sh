#!/bin/sh
current=`date "+%Y%m%d%H%M%S"`
docker logs facebook-spider-root > ./log/$current.log
docker-compose down
