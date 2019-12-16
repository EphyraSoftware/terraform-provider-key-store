package main

import (
	"github.com/ephyrasoftware/terraform-provider-keystore/impl"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"bundle": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func pkcsBundleCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	certPEM := d.Get("cert_pem").(string)
	keyPEM := d.Get("key_pem").(string)

	str, err := impl.CreateBundle(certPEM, keyPEM)
	if err != nil {
		return err
	}

	err = d.Set("bundle", str)
	if err != nil {
		return err
	}

	d.SetId(name)

	return pkcsBundleRead(d, m)
}

func pkcsBundleRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func pkcsBundleUpdate(d *schema.ResourceData, m interface{}) error {
	return pkcsBundleRead(d, m)
}

func pkcsBundleDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
