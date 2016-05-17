#!/bin/ash

curl http://eventsource:80/subscribe?channels=*,$1 -o "result_for_$1.txt"
