# Terraform Distro Image Provider

This [Terraform](http://terraform.io) provider is for dynamically finding distro images for a given IaaS.

## Status

Development/Testing

Currently only supports aws for coreos and ubuntu.

## Install - Build

```
go get
go build
go install
```

## With Homebrew

```
$ brew tap 'samsung-cnct/terraform-provider-coreosbox'
$ brew install terraform-provider-coreosbox
```

## Example Usage


```
data "distro_image" "coreos_ami" {
    cloud_provider = "aws"
    distribution = "coreos"
    region = "us-east-1"
    version = "current"
}

data "distro_image" "ubuntu_ami" {
    cloud_provider = "aws"
    distribution = "ubuntu"
    region = "us-west-2"
    store = "ssd"
    version = "16.04"
}

output "info_ami_coreos" {
    value = "Image: ${data.distro_image.coreos_ami.output_name}, ami: ${data.distro_image.coreos_ami.output_path}." 
}

output "info_ami_ubuntu" {
    value = "Image: ${data.distro_image.ubuntu_ami.output_name}, ami: ${data.distro_image.ubuntu_ami.output_path}." 
}
```

## Argument Reference

The following arguments are supported:

* arch - (optional) The type of architecture, ie. amd64 (default), x86_64, i686, etc. NOTE: This is distro specific.
* channel - (optional) CoreOS update channel, CoreOS only.
* cloud_provider - (required) The name of the cloud provider: aws, gke, jpc
* distribution - (optional) The distro to be used with the cloud provider: coreos (default), ubuntu
* region - (optional) The provider region in which to select the image from
* store - (optional) The provider store (if used), ie. instance, ebs, ssd
* subversion - (optional) The distro sub-verson to be used. Default to "latest".
* version - (optional) The distro version to be used. Default to "current".
* virtualization - (optional) The virtualization type. Provider specific, ie. pv or hvm for aws.

## Attributes Reference

The following attributes are exported:

* output_path - The path data needed to select this image, ami id for aws, etc.
* output_name - The human readable name of the image returned.

