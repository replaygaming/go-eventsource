# go-eventsource

Server-sent stream to update game info in real-time

## Build docker image

To build the `go-eventsource` docker image you should use [jet][1] by running
```
jet steps
```

This should generate an image tagged `us.gcr.io/replay-gaming/go-eventsource`
in your local docker engine.

Alternatively you can download the image from our [private docker
registry][2] manually too.

## Build locally

To build the binary outside of docker, you can use the standard commands from
golang:

```
go build
# generates ./go-eventsource
./go-eventsource -h
```

## Run

To run the binary, you can check all the options available by running
```
go-eventsource -h
```

All options can be overridden through environment variables following the
pattern `ES_<VARNAME>`. Eg. `ES_PORT=3333 ./go-eventsource` changes the http port
to `3333`

As always for more details on the options, read the source code :)


## Useful links

* https://golang.org/doc/install
* https://golang.org/doc/code.html

[1]: https://codeship.com/documentation/docker/installation/
[2]: https://replaygaming.atlassian.net/wiki/display/DT/Private+Docker+Registry
