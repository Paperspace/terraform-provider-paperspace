module github.com/Paperspace/terraform-paperspace-provider

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/hashicorp/terraform-plugin-sdk v1.13.1
	github.com/paperspace/paperspace-go v0.1.2
)

go 1.14

replace github.com/paperspace/paperspace-go => ../paperspace-go
