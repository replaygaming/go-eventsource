dependencies:
	go mod download

test:
	go test

compile: dependencies test
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -v -o bin/darwin_amd64
	GOOS=linux GOARCH=amd64 go build -v -o bin/linux_amd64
