# terraform-provider-paperspace
Paperspace terraform provider

## Release Notes

v0.1.2 -- first release; supports machine and script resources for create/read/destroy/import, and data sources for networks, templates, and users.

Note: currently this provider is offered as a terraform 'private cloud' provider while under early development.  We are moving toward contributing it back to the terraform open source project, which will remove the need for a separate download and installation step in the future.

## Downloads

Darwin (macOS) x64 [terraform-provider-paperspace-darwin](https://ps-terraform.s3.amazonaws.com/darwin/terraform-provider-paperspace)¹

Linux x64 [terraform-provider-paperspace](https://s3.amazonaws.com/ps-terraform/terraform-provider-paperspace)¹

Windows x64 [terraform-provider-paperspace.exe](https://s3.amazonaws.com/ps-terraform/terraform-provider-paperspace.exe)²

¹ Built with go 1.14 and tested with terraform 0.12.26
² Built with go 1.8.3 and tested with terraform 0.9.11

## Installation and Testing
1) Install [terraform](https://www.terraform.io/downloads.html) and make sure it is in your path.

2) Download the Paperspace terraform provider from one of the links above, or build it from source as described below.

3) Acquire a Paperspace API key for your account. See [Paperspace API](https://paperspace.github.io/paperspace-node/) for instructions on creating an api key.

4) A sample terraform config file is provided in [src/terraform-provider-paperspace/main.tf](src/terraform-provider-paperspace/main.tf)

Modify this file to use your actual api_key, and valid user email address in the account associated with the api_key.

5) cd to the directory where the sample configuration file is located, e.g.:
```
cd src/terraform-provider-paperspace
```

6) Make sure the `terraform-provider-paperspace` executable is also in your path, or in the directory with the .tf config files you want to use.

7) Run the following terraform commands interactively to exercise the configuration and examine the output:

(Note the sample configuration will create a machine with a public ip; testing this configuration will result in charges for the machine and public ip resources in most cases.)

```
terraform plan
terraform apply
terraform show
terraform refresh
terraform show
terraform plan
```

8) When you are done with testing, run the following to destroy the configuration (and thus destroy the machine and script objects created above):
```
terraform destroy
terraform show
```  

## Building from source

1) Install the latest version of [go](https://golang.org/dl/)

Note: this version of the provider has been successfully compiled with go version 1.8.3

2) Clone this repository and change to the project directory
```
git clone https://github.com/Paperspace/terraform-provider-paperspace.git
cd terraform-provider-paperspace
```

3) Change to the `src/terraform-provider-paperspace` subdirectory
```
cd src/terraform-provider-paperspace
```

4) Build the Paperspace terraform provider

On Mac run:
```
go mod tidy
go build ./...
```

On Linux x64 run:
```
make
```

On Windows run:
```
go mod tidy
go build
```

The output of the build  is a `terraform-provider-paperspace` executable.

Note: you cannot execute this provider binary directly.  The binary will be loaded by the terraform app if the provider binary is in your path and your .tf configuration files refer to the paperspace provider and paperspace resources, or datasources.

## Contributing

Want to contribute?  Contact us at hello@paperspace.com
