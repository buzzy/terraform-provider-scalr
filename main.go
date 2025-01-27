package main

import (
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/scalr/terraform-provider-scalr/scalr"
)

const (
	scalrProviderAddr = "registry.scalr.io/scalr/scalr"
)

func main() {
	var isDebug bool
	flag.BoolVar(&isDebug, "debug", false, "Start provider in debug mode.")
	flag.Parse()

	// Remove any date and time prefix in log package function output to
	// prevent duplicate timestamp and incorrect log level setting
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	schema.DescriptionKind = schema.StringMarkdown

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: scalr.Provider,
		ProviderAddr: scalrProviderAddr,
		Debug:        isDebug,
	})
}
