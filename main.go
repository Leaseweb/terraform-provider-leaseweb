package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
	"terraform-provider-leaseweb/internal/provider"
)

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name leaseweb

var (
	version = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/LeaseWeb/leaseweb",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.NewProvider(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
