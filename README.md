# go-eventsource
Server-sent stream to update game info in real-time

## Usage

Get the binary for your [distribution](https://github.com/replaygaming/go-eventsource/releases)

```shell
./linux_amd64 -h

  -amqp-queue string
        AMQP Queue name (default "eventsource")
  -amqp-url string
        AMQP URL (default "amqp://guest:guest@127.0.0.1:5672/eventsource")
  -compression
        Enable zlib compression of data
  -env string
        Environment: development or production (default "development")
  -port string
        Eventsource port (default "3001")
  -statsd-prefix string
        StatsD Prefix (default "app.es_go")
  -statsd-url string
        StatsD URL (default "127.0.0.1:8125")
```

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

## Contribuing

### Install `go`

Follow the instructions at [Golang.org](https://golang.org). **DO NOT** install using your distro pkg manager.

### Get project dependencies
    
    export GOPATH=~/go
    go get github.com/replaygaming/go-eventsource

### Running 

    cd ~/go/src/github.com/replaygaming/go-eventsource
    go build
    ./go-eventsource

### Important docs to read

* https://golang.org/doc/install
* https://golang.org/doc/code.html
