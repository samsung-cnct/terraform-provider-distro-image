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
		log.Printf("[DEBUG] Searching for CoreOS image for AWS.")
		version, err := getAwsAmiCoreOSVersion(d)
		if err != nil {
			return err
		}
		path, err := getAwsAmiCoreOSPath(d)
		if err != nil {
			return err
		}
		id := getAwsCoreOSId(d)
		setData(d, id, version, path)
		return nil
	case "ubuntu":
		log.Printf("[DEBUG] Searching for Ubuntu image for AWS.")
		version, err := getAwsAmiUbuntuVersion(d)
		if err != nil {
			return err
		}
		path, err := getAwsAmiUbuntuPath(d)
		if err != nil {
			return err
		}
		id := getAwsUbuntuId(d)
		setData(d, id, version, path)
		return nil
	}
	return fmt.Errorf("%s is not a supported AWS distribution\n", distro)
}

func setData(d *schema.ResourceData, id string, version string, path string) error {
	d.Set("output_name", fmt.Sprintf("%s-%s-%s", d.Get("distribution").(string), id, version))
	d.Set("output_path", path)
	d.SetId(id)
	return nil
}
