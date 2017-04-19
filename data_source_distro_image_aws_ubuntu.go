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
	Versions   map[string]Version `json:versions`
	Arch       string             `json:arch`
	Supported  bool               `json:supported`
	Version    string             `json:version`
	Aliases    string             `json:aliases`
	OS         string             `json:os`
	SupportEol string             `json:support_eol`
}

type Index struct {
	Product string
	Version string
	Item    string
}

type AwsUbuntuImages struct {
	resourceData *schema.ResourceData
	Products     map[string]Product
	index        Index
}

func NewAwsUbuntuImages(d *schema.ResourceData) (*AwsUbuntuImages, error) {
	i := new(AwsUbuntuImages)
	i.resourceData = d

	if err := i.getImages(); err != nil {
		return nil, err
	}

	if err := i.selectImage(); err != nil {
		return nil, err
	}
	return i, nil
}

func (i AwsUbuntuImages) GetVersion() string {
	return strings.Join([]string{i.Products[i.index.Product].Version, i.index.Version}, ".")
}

func (i AwsUbuntuImages) GetId() string {
	return strings.Join([]string{i.index.Product, i.index.Version, i.index.Item}, ":")
}

func (i AwsUbuntuImages) GetPath() string {
	return i.Products[i.index.Product].Versions[i.index.Version].Items[i.index.Item].Id
}

func (i *AwsUbuntuImages) getImages() error {
	type jsonData struct {
		Products map[string]Product `json:products`
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

	i.Products = data.Products

	return nil
}

func (i *AwsUbuntuImages) selectImage() error {
	arch := i.resourceData.Get("arch").(string)
	region := i.resourceData.Get("region").(string)
	store := i.resourceData.Get("store").(string)
	version := i.resourceData.Get("version").(string)
	subversion := i.resourceData.Get("subversion").(string)
	virtualization := i.resourceData.Get("virtualization").(string)
	maxversion := ""

	if !i.isValidArch(arch) {
		return fmt.Errorf("Invalid arch string, %s.\n", arch)
	}

	if !i.isValidStore(store) {
		return fmt.Errorf("Store, %s, is not available.\n", store)
	}

	if !i.isValidRegion(region) {
		return fmt.Errorf("Region, %s, is not available.\n", region)
	}

	if !i.isValidVirt(virtualization) {
		return fmt.Errorf("Virtualization method, %s, is not available.\n", virtualization)
	}

	// Walk the list of ubuntu product images looking for a match
	for productName, i_product := range i.Products {
		if i_product.Version != version || i_product.Arch != arch {
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
				if i_item.Crsn != region || i_item.RootStore != store || i_item.Virt != virtualization {
					continue
				}
				i.index.Item = itemName
				i.index.Product = productName
				i.index.Version = versionName
			}
		}
	}
	if i.index.Item != "" {
		return nil
	}

	return fmt.Errorf("No image match is found.")
}

func (i AwsUbuntuImages) isValidArch(arch string) bool {
	// Test to see if they offer an image for this arch
	for _, product := range i.Products {
		if arch == product.Arch {
			return true
		}
	}
	return false
}

func (i AwsUbuntuImages) isValidVirt(virt string) bool {
	// Test to see if they offer any image for this virtualization
	for _, product := range i.Products {
		for _, version := range product.Versions {
			for _, item := range version.Items {
				if virt == item.Virt {
					return true
				}
			}
		}
	}
	return false
}

func (i AwsUbuntuImages) isValidStore(store string) bool {
	// Test the image store to see if the requested store is valid, ie. instance, ebs, ssd, etc.
	for _, product := range i.Products {
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

func (i AwsUbuntuImages) isValidRegion(store string) bool {
	// Check the requested region against the list of image regions
	for _, product := range i.Products {
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
