![PS+Terraform2](https://user-images.githubusercontent.com/585865/90683337-5e825d00-e234-11ea-8bda-c4b299a00189.png)


# terraform-provider-paperspace
This is a Terraform provider for Paperspace infrastructure.

It is offered currently as a Terraform 'private cloud' provider while under early development.  We are moving toward contributing it back to the terraform open source project, which will remove the need for a separate download and installation step in the future.

[<img width="300" alt="Paperspace Terraform Provider " src="https://user-images.githubusercontent.com/585865/90682577-275f7c00-e233-11ea-9865-56e7f9205385.png">
](https://youtu.be/P3__yTs24rU)

[Read the blog post](https://blog.paperspace.com/introducing-paperspace-terraform-provider/)

## Installation and Testing
1) Install [terraform](https://www.terraform.io/downloads.html) v0.12.26 and make sure it is in your path.

2) Download the Paperspace terraform provider from [Releases](https://github.com/Paperspace/terraform-provider-paperspace/releases) (Linux, Windows, or Darwin), or [build it from source](#building-from-source).

3) Make terraform-provider-paperspace binary available as a plugin to your local Terraform by following the simple [Terraform plugin docs](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) – note: once you move it to your plugins directory, you may need to rename it to simply `terraform-provider-paperspace`

3) Acquire a Paperspace API key for your account. See [Paperspace API](https://paperspace.github.io/paperspace-node/) for instructions on creating an api key.

4) Copy the sample Terraform config file at [src/terraform-provider-paperspace/main.tf](src/terraform-provider-paperspace/main.tf) into your project directory.\
\
Modify this file to use your actual API Key, valid user email address, and team id for the account associated with the API Key.\
\
Note: if you clone down this repo, you can build/download the binary as a sibling to [src/terraform-provider-paperspace/main.tf](src/terraform-provider-paperspace/main.tf), replace the values described in #4 above with yours, and follow #5 below to use the Paperspace Terraform provider directly from this directory.

5) Run the following terraform commands interactively to exercise the configuration and examine the output.\
\
Note: the sample configuration will create a machine with a public ip; testing this configuration will result in charges for the machine and public ip resources in most cases.
    ```
    terraform plan
    terraform apply
    terraform show
    terraform refresh
    terraform show
    terraform plan
    ```

6) When you are done with testing, run the following to destroy the configuration (and thus destroy the machine and script objects created above):
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
