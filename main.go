package main

import (
	"flag"

	"github.com/basis-cloud/terraform-provider-bcc/bcc_terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool = true
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return bcc_terraform.Provider()
		},
	}

	if debugMode {
		opts.Debug = true
	}

	plugin.Serve(opts)
}
