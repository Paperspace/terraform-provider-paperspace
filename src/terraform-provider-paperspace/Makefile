build:
	go mod tidy
	go build -o terraform-provider-paperspace

build-linux:
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -o terraform-provider-paperspace-linux-amd64

build-darwin:
	go mod tidy
	GOOS=darwin GOARCH=amd64 go build -o terraform-provider-paperspace-darwin-amd64

build-windows:
	go mod tidy
	GOOS=windows GOARCH=amd64 go build -o terraform-provider-paperspace-windows-amd64.exe

build-all: build-linux build-darwin build-windows
