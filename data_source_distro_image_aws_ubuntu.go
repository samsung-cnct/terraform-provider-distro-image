package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type Item struct {
	RootStore string `json:"root_store"`
	Virt      string `json:virt`
	Crsn      string `json:crsn`
	Id        string `json:id`
}

type Version struct {
	Items   map[string]Item `json:items`
	Label   string          `json:label`
	PubName string          `json:pubname`
}

type Product struct {
	Versions        map[string]Version `json:versions`
	Arch            string             `json:arch`
	Supported       bool               `json:supported`
	ReleaseTitle    string             `json:release_title`
	ReleaseCodename string             `json:release_codename`
	Version         string             `json:version`
	Release         string             `json:release`
	Aliases         string             `json:aliases`
	OS              string             `json:os`
	SupportEol      string             `json:support_eol`
}

type Index struct {
	Product string
	Version string
	Item    string
}

var UbuntuImages map[string]Product

func getAwsAmiUbuntuVersion(d *schema.ResourceData) (string, error) {
	index, err := selectAwsUbuntuCloudImage(d)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{UbuntuImages[index.Product].Version, index.Version}, "."), nil
}

func getAwsUbuntuId(d *schema.ResourceData) string {
	index, err := selectAwsUbuntuCloudImage(d)
	if err != nil {
		return ""
	}

	return strings.Join([]string{index.Product, index.Version, index.Item}, ":")
}

func getAwsAmiUbuntuPath(d *schema.ResourceData) (string, error) {
	index, err := selectAwsUbuntuCloudImage(d)
	if err != nil {
		return "", err
	}
	return UbuntuImages[index.Product].Versions[index.Version].Items[index.Item].Id, nil
}

func getAwsAmiUbuntuCloudImages() error {
	// only fetch the data if it's not already cached.
	if len(UbuntuImages) > 0 {
		return nil
	}

	type jsonData struct {
		Updated    string             `json:updated`
		ApiVersion string             `json:format`
		DataType   string             `json:datatype`
		Products   map[string]Product `json:products`
		aliases    interface{}        `json:_aliases`
		ContentId  string             `json:content_id`
	}

	url := "https://cloud-images.ubuntu.com/releases/streams/v1/com.ubuntu.cloud:released:aws.json"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data jsonData
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return err
	}

	UbuntuImages = data.Products

	return nil
}

func selectAwsUbuntuCloudImage(d *schema.ResourceData) (Index, error) {
	var index Index

	err := getAwsAmiUbuntuCloudImages()
	if err != nil {
		return index, err
	}

	arch := d.Get("arch").(string)
	region := d.Get("region").(string)
	store := d.Get("store").(string)
	version := d.Get("version").(string)
	subversion := d.Get("subversion").(string)
	virtualization := d.Get("virtualization").(string)
	maxversion := ""

	if !isValidArch(arch) {
		return index, fmt.Errorf("Invalid arch string, %s.\n", arch)
	}

	if !isValidStore(store) {
		return index, fmt.Errorf("Store, %s, is not available.\n", store)
	}

	if !isValidRegion(region) {
		return index, fmt.Errorf("Region, %s, is not available.\n", region)
	}

	if !isValidVirt(virtualization) {
		return index, fmt.Errorf("Virtualization method, %s, is not available.\n", virtualization)
	}

	// Walk the list of ubuntu product images looking for a match
	for productName, i_product := range UbuntuImages {
		switch {
		case i_product.Version != version,
			i_product.Arch != arch:
			continue
		}
		for versionName, i_version := range i_product.Versions {
			if (subversion != "latest") && (versionName != subversion) {
				continue
			}
			if (subversion == "latest") && strings.Compare(versionName, maxversion) <= 0 {
				continue
			}
			maxversion = versionName
			for itemName, i_item := range i_version.Items {
				switch {
				case i_item.Crsn != region,
					i_item.RootStore != store,
					i_item.Virt != virtualization:
					continue
				}
				index.Item = itemName
				index.Product = productName
				index.Version = versionName
			}
		}
	}
	if index.Item != "" {
		return index, nil
	}

	return index, fmt.Errorf("No image match is found.")
}

func isValidArch(arch string) bool {
	switch arch {
	// These are the only two valid architectures for aws amis so far.
	case "amd64",
		"i386":
		return true
	}
	return false
}

func isValidVirt(virt string) bool {
	switch virt {
	// These are currently the only two virtualization types for aws amis.
	case "hvm",
		"pv":
		return true
	}
	return false
}

func isValidStore(store string) bool {
	// Test the image store to see if the requested store is valid, ie. instance, ebs, ssd, etc.
	for _, product := range UbuntuImages {
		for _, version := range product.Versions {
			for _, item := range version.Items {
				if store == item.RootStore {
					return true
				}
			}
		}
	}
	return false
}

func isValidRegion(store string) bool {
	// Check the requested region against the list of image regions
	for _, product := range UbuntuImages {
		for _, version := range product.Versions {
			for _, item := range version.Items {
				if store == item.Crsn {
					return true
				}
			}
		}
	}
	return false
}
