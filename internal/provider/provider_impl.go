package provider

import (
	"context"
	"os"
	"regexp"

	"github.com/eaglespirittech/terraform-provider-controltowermanagement/internal/client"
	"github.com/eaglespirittech/terraform-provider-controltowermanagement/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &controltowermanagementProvider{}
)

func (p *controltowermanagementProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "controltowermanagement"
	resp.Version = p.version
}

func (p *controltowermanagementProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Control Tower Management API.",
		Blocks: map[string]schema.Block{
			"assume_role": schema.SingleNestedBlock{
				Description: "Assume role configuration block",
				Attributes: map[string]schema.Attribute{
					"role_arn": schema.StringAttribute{
						Description: "ARN of the role to assume",
						Required:    true,
						Validators: []validator.String{
							validators.RegexMatches(
								regexp.MustCompile(`^arn:aws:iam::\d{12}:role/[a-zA-Z0-9_+=,.@-]+$`),
								"Role ARN must be a valid AWS IAM role ARN",
							),
						},
					},
					"session_name": schema.StringAttribute{
						Description: "Session name to use when assuming the role",
						Optional:    true,
					},
					"external_id": schema.StringAttribute{
						Description: "External ID to use when assuming the role",
						Optional:    true,
					},
					"duration_seconds": schema.Int64Attribute{
						Description: "Duration of the assumed role session in seconds",
						Optional:    true,
					},
					"policy_arns": schema.ListAttribute{
						Description: "List of ARNs of IAM policies to use for the assumed role session",
						ElementType: types.StringType,
						Optional:    true,
					},
					"policy": schema.StringAttribute{
						Description: "IAM policy document to use for the assumed role session",
						Optional:    true,
					},
					"tags": schema.MapAttribute{
						Description: "Map of tags to use for the assumed role session",
						ElementType: types.StringType,
						Optional:    true,
					},
					"transitive_tag_keys": schema.ListAttribute{
						Description: "List of tag keys to pass to the assumed role session",
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"access_key": schema.StringAttribute{
				Description: "AWS Access Key ID. Can be set via AWS_ACCESS_KEY environment variable.",
				Optional:    true,
			},
			"secret_key": schema.StringAttribute{
				Description: "AWS Secret Access Key. Can be set via AWS_SECRET_ACCESS_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"region": schema.StringAttribute{
				Description: "AWS Region. Can be set via AWS_REGION environment variable.",
				Optional:    true,
				Validators: []validator.String{
					validators.RegexMatches(
						regexp.MustCompile(`^[a-z]{2}-[a-z]+-\d{1}$`),
						"Region must be a valid AWS region (e.g., us-west-2)",
					),
				},
			},
		},
	}
}

func (p *controltowermanagementProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config controltowermanagementProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get credentials from environment variables if not set in configuration
	if config.AccessKey.IsNull() {
		accessKey := os.Getenv("AWS_ACCESS_KEY")
		if accessKey == "" {
			accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		}
		config.AccessKey = types.StringValue(accessKey)
	}
	if config.SecretKey.IsNull() {
		config.SecretKey = types.StringValue(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	}
	if config.Region.IsNull() {
		config.Region = types.StringValue(os.Getenv("AWS_REGION"))
	}

	// Validate required credentials
	if config.AccessKey.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing AWS Access Key",
			"Either set AWS_ACCESS_KEY/AWS_ACCESS_KEY_ID environment variable or provide access_key in provider configuration",
		)
	}
	if config.SecretKey.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing AWS Secret Key",
			"Either set AWS_SECRET_ACCESS_KEY environment variable or provide secret_key in provider configuration",
		)
	}
	if config.Region.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing AWS Region",
			"Either set AWS_REGION environment variable or provide region in provider configuration",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize AWS client with session token if available
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")
	awsClient, err := client.NewClient(
		config.AccessKey.ValueString(),
		config.SecretKey.ValueString(),
		config.Region.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create AWS client",
			"Error creating AWS client: "+err.Error(),
		)
		return
	}

	// Set session token if available
	if sessionToken != "" {
		if err := awsClient.SetSessionToken(sessionToken); err != nil {
			resp.Diagnostics.AddError(
				"Failed to set session token",
				"Error setting session token: "+err.Error(),
			)
			return
		}
	}

	// Handle assume role if configured
	if config.AssumeRole != nil {
		assumeRoleConfig := &client.AssumeRoleConfig{
			RoleArn: config.AssumeRole.RoleArn.ValueString(),
		}

		if !config.AssumeRole.SessionName.IsNull() {
			assumeRoleConfig.SessionName = config.AssumeRole.SessionName.ValueString()
		}
		if !config.AssumeRole.ExternalId.IsNull() {
			assumeRoleConfig.ExternalId = config.AssumeRole.ExternalId.ValueString()
		}
		if !config.AssumeRole.DurationSeconds.IsNull() {
			assumeRoleConfig.DurationSeconds = int32(config.AssumeRole.DurationSeconds.ValueInt64())
		}
		if !config.AssumeRole.Policy.IsNull() {
			assumeRoleConfig.Policy = config.AssumeRole.Policy.ValueString()
		}

		// Handle policy ARNs
		if !config.AssumeRole.PolicyArns.IsNull() {
			var policyArns []string
			diags := config.AssumeRole.PolicyArns.ElementsAs(ctx, &policyArns, false)
			if !diags.HasError() {
				assumeRoleConfig.PolicyArns = policyArns
			}
		}

		// Handle tags
		if !config.AssumeRole.Tags.IsNull() {
			var tags map[string]string
			diags := config.AssumeRole.Tags.ElementsAs(ctx, &tags, false)
			if !diags.HasError() {
				assumeRoleConfig.Tags = tags
			}
		}

		// Handle transitive tag keys
		if !config.AssumeRole.TransitiveTagKeys.IsNull() {
			var tagKeys []string
			diags := config.AssumeRole.TransitiveTagKeys.ElementsAs(ctx, &tagKeys, false)
			if !diags.HasError() {
				assumeRoleConfig.TransitiveTagKeys = tagKeys
			}
		}

		// Assume the role
		if err := awsClient.AssumeRole(ctx, assumeRoleConfig); err != nil {
			resp.Diagnostics.AddError(
				"Failed to assume role",
				"Error assuming role: "+err.Error(),
			)
			return
		}
	}

	// Make the AWS client available during DataSource and Resource type Configure methods
	resp.DataSourceData = awsClient
	resp.ResourceData = awsClient
}

func (p *controltowermanagementProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAwsAccountDataSource,
	}
}

func (p *controltowermanagementProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
