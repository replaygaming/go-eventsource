ROOT := $(shell pwd)
default:
	CGO_CFLAGS="-I$(ROOT)/bin/include" \
	CGO_LDFLAGS="-L$(ROOT)/bin/lib -lnewrelic-collector-client -lnewrelic-common -lnewrelic-transaction" \
	GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=amd64 go build -v -o bin/eventsource
