package main

import (
	"encoding/base64"
	"github.com/ephyrasoftware/terraform-provider-keystore/impl"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"os"
	"path"
)

func pkcsBundle() *schema.Resource {
	return &schema.Resource{
		Create: pkcsBundleCreate,
		Read:   pkcsBundleRead,
		Update: pkcsBundleUpdate,
		Delete: pkcsBundleDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cert_pem": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"key_pem": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ca_certs": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"bundle": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func pkcsBundleCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	outputPath := m.(KeystoreConfig).Path

	certPEM := d.Get("cert_pem").(string)
	keyPEM := d.Get("key_pem").(string)
	caCerts := d.Get("ca_certs").(*schema.Set).List()

	err := impl.CreateBundle(certPEM, keyPEM, impl.SliceOfString(caCerts), outputPath, name)
	if err != nil {
		return err
	}

	d.SetId(name)

	return pkcsBundleRead(d, m)
}

func pkcsBundleRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	outputPath := m.(KeystoreConfig).Path

	var outFile = path.Join(outputPath, name+".p12")

	bundle, err := ioutil.ReadFile(outFile)
	if err != nil {
		d.SetId("")
		return nil
	}

	err = d.Set("bundle", base64.StdEncoding.EncodeToString(bundle))
	if err != nil {
		return err
	}

	return nil
}

func pkcsBundleUpdate(d *schema.ResourceData, m interface{}) error {
	return pkcsBundleRead(d, m)
}

func pkcsBundleDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	outputPath := m.(KeystoreConfig).Path

	var outFile = path.Join(outputPath, name+".p12")

	err := os.Remove(outFile)

	return err
}
