#!/bin/ash

curl -N http://eventsource:80/subscribe?channels=*,$1 | tee "result_for_${1}.txt" > /dev/null
