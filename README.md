# go-eventsource
Server-sent stream to update game info in real-time

## Configure RabbitMQ

### Install `rabbitmq` and `rabbitmqadmin`

Download and installation guide from [RabbitMQ site](https://www.rabbitmq.com/download.html).
rabbitmqadmin is binary, found as part of [rabbitmq-management](https://github.com/rabbitmq/rabbitmq-management) project.

### Enable the management plugin:

    [sudo] rabbitmq-plugins enable rabbitmq_management

Then (re)start the rabbitmq daemon.

    [sudo] sudo rabbitmqctl stop
    [sudo] rabbitmq-server -detached

Declare the host and exchange for the eventsource:

    rabbitmqadmin declare vhost name=eventsource
    rabbitmqadmin declare permission vhost=eventsource user=guest configure=".*" write=".*" read=".*"
    rabbitmqadmin -V eventsource declare exchange name=es_ex type=fanout durable=true

## Configure Go

### Install `golang`

    [sudo] apt-get install golang

Set GOPATH

    # cd to project root
    export GOPATH=`pwd`

### Get project dependencies

    go get .

### Cross-compiling

    make compile

## Run `go-eventsource`

OS: Linux ARCH: amd64

    ./bin/linux_amd64

OS: Darwin ARCH amd64

    ./bin/darwin_amd64

You can change `GOOS` and `GOARCH` in `Makefile`.
