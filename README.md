# terraform-provider-paperspace
This is an Terraform provider for Paperspace infrastructure.

It is offered currently as a Terraform 'private cloud' provider while under early development.  We are moving toward contributing it back to the terraform open source project, which will remove the need for a separate download and installation step in the future.


## Releases
Visit [Releases](https://github.com/Paperspace/terraform-provider-paperspace/releases) to download a pre-compiled binary for Linux, Windows, or Darwin. You can also [compile from source](#building-from-source).


## Installation and Testing
1) Install [terraform](https://www.terraform.io/downloads.html) v0.12.26 and make sure it is in your path.

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

1) Install the latest version of [go](https://golang.org/dl/) that supports go modules (we currently use go 1.14 for this project)

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
On any platform:
```
make build
```

5) Compile the provider binary for various platforms

For Linux x64:
```
make build-linux
```

For Windows x64:
```
make build-windows
```

For Darwin x64 (macOS):
```
make build-darwin
```

The output of the build is a `terraform-provider-paperspace` executable.

Note: you cannot execute this provider binary directly.  The binary will be loaded by the terraform app if the provider binary is in your path and your .tf configuration files refer to the paperspace provider and paperspace resources, or datasources.

## Contributing

Want to contribute? Contact us at hello@paperspace.com
