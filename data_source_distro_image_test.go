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
						regexp.MustCompile(`^CoreOS-stable-[0-9.]+$`),
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
