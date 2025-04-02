package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function is invoked for every Terraform CLI
// command executed to create a provider server to which the CLI can reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"controltowermanagement": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example
	// assertions about the environment to help prevent test failures in CI/CD.
	// For example:
	// os.Setenv("AWS_PROFILE", "sandbox")
	// os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
}

func TestAccAwsAccountDataSource(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsAccountDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.controltowermanagement_aws_account.test", "accounts.#"),
					resource.TestCheckResourceAttrSet("data.controltowermanagement_aws_account.test", "accounts.0.account_id"),
					resource.TestCheckResourceAttrSet("data.controltowermanagement_aws_account.test", "accounts.0.account_name"),
					resource.TestCheckResourceAttrSet("data.controltowermanagement_aws_account.test", "accounts.0.email"),
					resource.TestCheckResourceAttrSet("data.controltowermanagement_aws_account.test", "accounts.0.status"),
				),
			},
		},
	})
}

func testAccAwsAccountDataSourceConfig() string {
	return `
provider "controltowermanagement" {
  access_key = "` + os.Getenv("AWS_ACCESS_KEY") + `"
  secret_key = "` + os.Getenv("AWS_SECRET_ACCESS_KEY") + `"
  region     = "` + os.Getenv("AWS_REGION") + `"
}

data "controltowermanagement_aws_account" "test" {}
`
}
