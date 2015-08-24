# go-eventsource
Server-sent stream to update game info in real-time

## Configure RabbitMQ

Install rabbitmq and the rabbitmqadmin (both available on AUR for Arch Linux) and enable the management plugin:

    [sudo] rabbitmq-plugins enable rabbitmq_management

then (re)start the rabbitmq daemon. Declare the host and exchange for the eventsource:

    rabbitmqadmin declare vhost name=eventsource

    rabbitmqadmin declare permission vhost=eventsource user=guest configure=".*" write=".*" read=".*"

    rabbitmqadmin -V eventsource declare exchange name=es_ex type=fanout durable=true

## Cross-compiling

    make compile
    ls bin

