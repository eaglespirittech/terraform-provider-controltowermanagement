package client

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrganizationsAPI is a mock implementation of the Organizations API
type MockOrganizationsAPI struct {
	mock.Mock
}

func (m *MockOrganizationsAPI) ListAccounts(ctx context.Context, params *organizations.ListAccountsInput, optFns ...func(*organizations.Options)) (*organizations.ListAccountsOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*organizations.ListAccountsOutput), args.Error(1)
}

func TestGetAccountInfo(t *testing.T) {
	mockOrg := new(MockOrganizationsAPI)

	// Test data
	testAccounts := []types.Account{
		{
			Id:     aws.String("123456789012"),
			Name:   aws.String("Test Account 1"),
			Email:  aws.String("test1@example.com"),
			Status: types.AccountStatusActive,
		},
		{
			Id:     aws.String("098765432109"),
			Name:   aws.String("Test Account 2"),
			Email:  aws.String("test2@example.com"),
			Status: types.AccountStatusActive,
		},
	}

	// Set up expectations
	mockOrg.On("ListAccounts", mock.Anything, &organizations.ListAccountsInput{}).
		Return(&organizations.ListAccountsOutput{
			Accounts: testAccounts,
		}, nil)

	// Create test client
	testClient := &Client{
		awsConfig: aws.Config{},
	}

	// Execute test
	accounts, err := testClient.GetAccountInfo(context.Background())

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, accounts, 2)
	assert.Equal(t, "123456789012", accounts[0].AccountId)
	assert.Equal(t, "Test Account 1", accounts[0].AccountName)
	assert.Equal(t, "test1@example.com", accounts[0].Email)
	assert.Equal(t, "ACTIVE", accounts[0].Status)
}

func TestAssumeRole(t *testing.T) {
	// Create test client
	testClient := &Client{
		awsConfig: aws.Config{},
	}

	// Test assume role configuration
	assumeRoleConfig := &AssumeRoleConfig{
		RoleArn:         "arn:aws:iam::123456789012:role/TestRole",
		SessionName:     "test-session",
		DurationSeconds: 3600,
	}

	// Execute test
	err := testClient.AssumeRole(context.Background(), assumeRoleConfig)

	// Note: This is a basic test that just verifies the function doesn't panic
	// In a real test environment, you would mock the STS client and verify the credentials
	assert.NoError(t, err)
}
