package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDistroImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDistroImageRead,

		Schema: map[string]*schema.Schema{
			"cloud_provider": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUD_PROVIDER", nil),
				Description: "A cloud provider name: aws, gke, jpc, etc.",
				ForceNew:    true,
			},
			"distribution": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUD_DISTRO", "coreos"),
				Description: "The distro to be used with the cloud provider: coreos (default), arch, ubuntu, etc",
				ForceNew:    true,
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUD_DISTRO_VERSION", "latest"),
				Description: "The distro version to be used. Default to \"latest\".",
				ForceNew:    true,
			},
			"output_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "The image name to use for this cloud provider.",
			},
		},
	}
}

func dataSourceDistroImageRead(d *schema.ResourceData, meta interface{}) error {
	name := "xyz_image.ami"
	d.Set("output_name", name)
	log.Printf("%k", d)
	d.SetId(name)
	return nil
}
