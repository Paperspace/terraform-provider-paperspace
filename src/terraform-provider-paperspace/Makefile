build-linux:
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -o terraform-provider-paperspace

build-darwin:
	go mod tidy
	GOOS=darwin GOARCH=amd64 go build -o terraform-provider-paperspace
