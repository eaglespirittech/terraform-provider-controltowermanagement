package provider

import (
	"context"
	"fmt"

	"github.com/eaglespirittech/terraform-provider-controltowermanagement/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &awsAccountDataSource{}
	_ datasource.DataSourceWithConfigure = &awsAccountDataSource{}
)

// awsAccountDataSource is the data source implementation.
type awsAccountDataSource struct {
	client *client.Client
}

// awsAccountDataSourceModel describes the data source data model.
type awsAccountDataSourceModel struct {
	Accounts []awsAccountModel `tfsdk:"accounts"`
}

// awsAccountModel describes the AWS account model.
type awsAccountModel struct {
	AccountId   types.String `tfsdk:"account_id"`
	AccountName types.String `tfsdk:"account_name"`
	Email       types.String `tfsdk:"email"`
	Status      types.String `tfsdk:"status"`
}

// NewAwsAccountDataSource is a helper function to simplify the provider implementation.
func NewAwsAccountDataSource() datasource.DataSource {
	return &awsAccountDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *awsAccountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *awsAccountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_account"
}

// Schema defines the schema for the data source.
func (d *awsAccountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about AWS accounts in your organization.",
		Attributes: map[string]schema.Attribute{
			"accounts": schema.ListNestedAttribute{
				Description: "List of AWS accounts in the organization",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Description: "The ID of the AWS account",
							Computed:    true,
						},
						"account_name": schema.StringAttribute{
							Description: "The name of the AWS account",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "The email address associated with the account",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The status of the account (ACTIVE, SUSPENDED, etc.)",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *awsAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state awsAccountDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Configured",
			"Expected configured client but got nil. Please report this issue to the provider developers.",
		)
		return
	}

	accounts, err := d.client.GetAccountInfo(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading AWS Accounts",
			fmt.Sprintf("Could not read AWS accounts: %s\n\nThis could be due to:\n1. Invalid AWS credentials\n2. Missing required IAM permissions (organizations:ListAccounts)\n3. Invalid AWS region configuration\n4. Network connectivity issues\n\nPlease check your AWS credentials and permissions.", err.Error()),
		)
		return
	}

	// Map response body to model
	for _, account := range accounts {
		accountState := awsAccountModel{
			AccountId:   types.StringValue(account.AccountId),
			AccountName: types.StringValue(account.AccountName),
			Email:       types.StringValue(account.Email),
			Status:      types.StringValue(account.Status),
		}
		state.Accounts = append(state.Accounts, accountState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
