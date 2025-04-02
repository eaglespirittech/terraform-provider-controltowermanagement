package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// controltowermanagementProvider is the provider implementation.
type controltowermanagementProvider struct {
	version string
}

// assumeRoleModel represents the assume role configuration
type assumeRoleModel struct {
	RoleArn           types.String `tfsdk:"role_arn"`
	SessionName       types.String `tfsdk:"session_name"`
	ExternalId        types.String `tfsdk:"external_id"`
	DurationSeconds   types.Int64  `tfsdk:"duration_seconds"`
	PolicyArns        types.List   `tfsdk:"policy_arns"`
	Policy            types.String `tfsdk:"policy"`
	Tags              types.Map    `tfsdk:"tags"`
	TransitiveTagKeys types.List   `tfsdk:"transitive_tag_keys"`
}

// controltowermanagementProviderModel describes the provider data model.
type controltowermanagementProviderModel struct {
	AccessKey  types.String     `tfsdk:"access_key"`
	SecretKey  types.String     `tfsdk:"secret_key"`
	Region     types.String     `tfsdk:"region"`
	AssumeRole *assumeRoleModel `tfsdk:"assume_role"`
}
