package main

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
)

func TestAccDistroImage_Basic(t *testing.T) {
	r.Test(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			r.TestStep{
				Config: testAccDistroImageAwsCoreOSConfig,
				Check: r.ComposeTestCheckFunc(
					r.TestMatchResourceAttr(
						"data.distro_image.foo", "output_name",
						regexp.MustCompile(`^coreos-stable:current-[0-9.]+$`),
					),
				),
			},
			r.TestStep{
				Config: testAccDistroImageAwsUbuntuConfig,
				Check: r.ComposeTestCheckFunc(
					r.TestMatchResourceAttr(
						"data.distro_image.bar", "output_name",
						regexp.MustCompile(`^ubuntu-com.ubuntu.cloud:server:16.04:amd64:[0-9.]+:[a-z0-9]+-[0-9.]+$`),
					),
				),
			},
		},
	})
}

var testAccDistroImageAwsCoreOSConfig = `
data "distro_image" "foo" {
    cloud_provider = "aws"
}
`

var testAccDistroImageAwsUbuntuConfig = `
data "distro_image" "bar" {
    cloud_provider = "aws"
	distribution = "ubuntu"
	store = "ssd"
	version = "16.04"
}
`
