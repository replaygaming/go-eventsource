FROM golang:1.13-alpine AS build_base

WORKDIR /src

COPY go.sum go.mod ./

# COPY . .
COPY *.go ./

# Unit tests
RUN go build -o /out .

# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add -U --no-cache ca-certificates

COPY --from=build_base /out /app/go-eventsource

ENTRYPOINT ["/app/go-eventsource"]