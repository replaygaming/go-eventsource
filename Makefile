compile: dependencies
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -v -o bin/darwin_amd64
	GOOS=linux GOARCH=amd64 go build -v -o bin/linux_amd64

dependencies:
	go get
