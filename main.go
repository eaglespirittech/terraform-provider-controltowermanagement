package main

import (
	"context"

	"github.com/eaglespirittech/terraform-provider-controltowermanagement/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/eaglespirittech/controltowermanagement",
	}
	providerserver.Serve(context.Background(), provider.New("dev"), opts)
}
