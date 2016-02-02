# go-eventsource
Server-sent stream to update game info in real-time

## Usage
Get the binary for your [distribution](https://github.com/replaygaming/go-eventsource/releases)

Or build it yourself:

Get project dependencies

    export GOPATH=~/go
    go get github.com/replaygaming/go-eventsource
    cd ~/go/src/github.com/replaygaming/go-eventsource
    go get

Compile and run

    go build
    ./go-eventsource

### Configuration
Configuration is done using environment variables

```shell
# Environment
ENV=[development|production]

# Eventsource port (for example "3001")
PORT

# AMQP URL (for example "amqp://guest:guest@127.0.0.1:5672/eventsource")
AMQP_URL

# AMQP Queue name (for example "eventsource")
AMQP_QUEUE

# StatsD URL (for example "127.0.0.1:8125")
STATSD_URL

# StatsD Prefix (for example "app.es_go")
STATSD_PREFIX

# Enable zlib compression of data
COMPRESS=[true|false]
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

## Important docs to read

* https://golang.org/doc/install
* https://golang.org/doc/code.html
