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

### Install

#### Go on OS X

Using Homebrew: `brew install go` or with support for crossc ompilation `brew install go --with-cc-common`

#### Go on Linux

Follow the instructions at [Golang.org](https://golang.org). **DO NOT** install using your distro pkg manager.

#### Get code
  
Assuming you have some directory, `~/projects` in this example, where you want to keep the source:

    cd ~/projects
    git clone git@github.com:replaygaming/go-eventsource.git

#### Go environment
    
    export GOPATH=$GOPATH:~/projects/go-eventsource

#### Compiling and running

    cd ~/projects/go-eventsource/src/eventsource-server
    go get .
    go build
    ./eventsource-server

### Important docs to read

* https://golang.org/doc/install
* https://golang.org/doc/code.html
