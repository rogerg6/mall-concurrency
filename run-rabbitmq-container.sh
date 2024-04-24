#!/bin/bash

RABBITMQ_CONTAINER=rabbitmq

num=`docker ps -a | grep -w "$RABBITMQ_CONTAINER" | wc -l`

if [ 0 -eq $num ]; then
  docker run -it  -p 5672:5672 -p 15672:15672 --name rabbitmq rabbitmq:3.13-management /bin/bash
else
  docker start $RABBITMQ_CONTAINER
  docker exec -it -w /root $RABBITMQ_CONTAINER /bin/bash
fi
