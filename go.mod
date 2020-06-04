module github.com/Paperspace/terraform-paperspace-provider

go 1.14

require (
	github.com/aws/aws-sdk-go v1.30.12 // indirect
	github.com/go-resty/resty v1.12.0
	github.com/hashicorp/go-getter v1.4.2-0.20200106182914-9813cbd4eb02 // indirect
	github.com/hashicorp/go-plugin v1.3.0 // indirect
	github.com/hashicorp/hcl/v2 v2.3.0 // indirect
	github.com/hashicorp/terraform-config-inspect v0.0.0-20191212124732-c6ae6269b9d7 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.13.1
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37 // indirect
)

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.12.0
