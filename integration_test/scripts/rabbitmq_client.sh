#!/bin/ash

./rabbitmqadmin -H rabbitmq -u guest -p guest publish exchange=es_ex routing_key="all" payload="{\"data\":\"$1\"}"
