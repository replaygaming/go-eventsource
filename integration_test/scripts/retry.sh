#!/usr/bin/env ash

sleep 20

count=0

until ( /go-eventsource )
do
  count=$((count + 1))
  echo "Starting go-eventsource failed, retrying in 1 sec."
  if [ ${count} -gt 10 ]
  then
    echo "Services didn't become ready in time"
    exit 1
  fi
  sleep 1
done
