package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDistroImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDistroImageRead,

		Schema: map[string]*schema.Schema{
			"arch": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Archetecture string, ie. amd64, x86_64, i686, etc.",
				Default:     "amd64",
				Optional:    true,
				ForceNew:    true,
			},
			"channel": &schema.Schema{
				Type:        schema.TypeString,
				Description: "CoreOS update channel",
				Default:     "stable",
				Optional:    true,
				ForceNew:    true,
			},
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
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Region. Applies to ami boxes.",
				Default:     "us-west-2",
				Optional:    true,
				ForceNew:    true,
			},
			"store": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Storage Type.",
				Default:     "ebs",
				Optional:    true,
				ForceNew:    true,
			},
			"subversion": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUD_DISTRO_SUBVERSION", "latest"),
				Description: "The distro version to be used. Default to \"latest\".",
				ForceNew:    true,
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUD_DISTRO_VERSION", "current"),
				Description: "The distro version to be used. Default to \"current\".",
				ForceNew:    true,
			},
			"virtualization": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Virtualization type. Applies to ami boxes: pv or hvm",
				Default:     "hvm",
				Optional:    true,
				ForceNew:    true,
			},
			"output_path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "Path information string, ie. gce path, vagrant box url, etc.",
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
	cloud_provider := d.Get("cloud_provider")
	switch cloud_provider {
	case "aws":
		return dataSourceDistroImageAwsRead(d, meta)
	case "jpc":
		return dataSourceDistroImageJpcRead(d, meta)
	}
	return fmt.Errorf("%s is not a supported cloud provider", cloud_provider)
}
