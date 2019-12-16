package main

import (
	"errors"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KEYSTORE_PATH", ""),
				Description: "The path to use for storing bundles.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"keystore_pkcs12_bundle": dataSourceKeystorePkcs12Bundle(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"keystore_pkcs12_bundle": pkcsBundle(),
		},

		ConfigureFunc: configureProvider,
	}
}

type KeystoreConfig struct {
	Path string
}

func configureProvider(d *schema.ResourceData) (meta interface{}, err error) {
	path := d.Get("path").(string)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, errors.New("path provided to keystore provider does not exist and cannot be created")
		}
	}

	return KeystoreConfig{Path: path}, nil
}
