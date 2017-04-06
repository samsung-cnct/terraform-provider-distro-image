package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDistroImageJpcRead(d *schema.ResourceData, meta interface{}) error {
	distro := d.Get("distribution")
	name := ""
	d.Set("output_name", name)
	d.SetId(name)
	return fmt.Errorf("%s is not a supported distribution", distro)
}
