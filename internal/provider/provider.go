package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// New creates a new provider instance
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &controltowermanagementProvider{
			version: version,
		}
	}
}
