package main

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataSourceKeystorePkcs12Bundle() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeystorePkcs12BundleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the bundle",
			},
			"bundle": {
				Type:        schema.TypeString,
				Description: "The PKCS12 bundle",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func dataSourceKeystorePkcs12BundleRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("name").(string))

	return pkcsBundleRead(d, meta)
}
