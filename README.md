# go-eventsource

Server-sent stream to update game info in real-time

## Build docker image

```bash
docker build -t us.gcr.io/replay-gaming/go-eventsource .
```

This should generate an image tagged `us.gcr.io/replay-gaming/go-eventsource:latest`
in your local docker engine.

Alternatively you can download the image from our [private docker
registry][1] manually too.

## Build locally

To build the binary outside of docker, you can use the standard commands from
golang:

```bash
go build
# generates ./go-eventsource
./go-eventsource -h
```

## CI

You always can monitor CI status in [cloudbuild]

## Run

To run the binary, you can check all the options available by running

```bash
go-eventsource -h
```

All options can be overridden through environment variables following the
pattern `ES_<VARNAME>`. Eg. `ES_PORT=3333 ./go-eventsource` changes the http port
to `3333`

As always for more details on the options, read the source code :)

## Dependencies

For tests you will need pubsub-emulator. This can be set with environment
variable: `PUBSUB_EMULATOR_HOST`, default value is `pubsub-emulator:8538`

## Provision

Use Helm for provision, example:

```bash
helm install charts/go-eventsource
```

## Useful links

* https://golang.org/doc/install
* https://golang.org/doc/code.html

[1]: https://replaygaming.atlassian.net/wiki/display/DT/Private+Docker+Registry
[cloudbuild]: https://console.cloud.google.com/cloud-build/builds?project=replay-gaming