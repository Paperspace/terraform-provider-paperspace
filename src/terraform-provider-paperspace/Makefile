build:
	go mod tidy
	go build -o terraform-provider-paperspace

build-linux:
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -o terraform-provider-paperspace

build-darwin:
	go mod tidy
	GOOS=darwin GOARCH=amd64 go build -o terraform-provider-paperspace

build-windows:
	go mod tidy
	GOOS=windows GOARCH=amd64 go build -o terraform-provider-paperspace.exe
