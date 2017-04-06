package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDistroImageAwsRead(d *schema.ResourceData, meta interface{}) error {
	distro := d.Get("distribution")
	switch distro {
	case "coreos":
		log.Printf("[DEBUG] Searching aws for CoreOS image.")
		version, err := getAwsAmiCoreOSVersion(d)
		if err != nil {
			return err
		}
		path, err := getAwsAmiCoreOSPath(d)
		if err != nil {
			return err
		}
		d.Set("output_name", version)
		d.Set("output_path", path)
		d.SetId(getAwsCoreOSId(d))
		return nil
	}
	return fmt.Errorf("%s is not a supported AWS distribution\n", distro)
}
