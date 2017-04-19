package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dlintw/goconf"
	"github.com/hashicorp/terraform/helper/schema"
)

func getAwsAmiCoreOSVersion(d *schema.ResourceData) (string, error) {
	url := fmt.Sprintf(
		"http://%s.release.core-os.net/amd64-usr/%s/version.txt",
		d.Get("channel").(string),
		d.Get("version").(string))

	log.Printf("[DEBUG] CoreOS lookup url: %v", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	verInfo, err := goconf.ReadConfigBytes(bodyBytes)
	if err != nil {
		return "", err
	}

	return verInfo.GetString("default", "COREOS_VERSION")
}

func getAwsCoreOSId(d *schema.ResourceData) string {
	channel := d.Get("channel").(string)
	v := d.Get("version").(string)

	return strings.Join([]string{channel, v}, ":")
}

func getAwsAmiCoreOSPath(d *schema.ResourceData) (string, error) {
	type (
		ami struct {
			Name string `json:"name"`
			PV   string `json:"pv"`
			HVM  string `json:"hvm"`
		}

		amiInfo struct {
			AMIs []ami `json:"amis"`
		}
	)

	url := fmt.Sprintf(
		"http://%s.release.core-os.net/amd64-usr/%s/coreos_production_ami_all.json",
		d.Get("channel").(string),
		d.Get("version").(string))

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data amiInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	virt := d.Get("virtualization").(string)
	region := d.Get("region").(string)

	for _, a := range data.AMIs {
		if a.Name == region {
			switch virt {
			case "pv":
				return a.PV, nil
			case "hvm":
				return a.HVM, nil
			default:
				return "", fmt.Errorf("Unknown virtualization type")
			}
		}
	}
	return "", fmt.Errorf("No ami found")
}
