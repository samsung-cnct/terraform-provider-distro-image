package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"distro_image": dataSourceDistroImage(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"distro_image": schema.DataSourceResourceShim(
				"distro_image",
				dataSourceDistroImage(),
			),
		},
	}
}
