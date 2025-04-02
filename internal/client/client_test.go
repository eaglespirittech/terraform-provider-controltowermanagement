package client

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgTypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	stsTypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
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

type mockOrganizationsClient struct {
	ListAccountsFunc func(ctx context.Context, params *organizations.ListAccountsInput, optFns ...func(*organizations.Options)) (*organizations.ListAccountsOutput, error)
}

func (m *mockOrganizationsClient) ListAccounts(ctx context.Context, params *organizations.ListAccountsInput, optFns ...func(*organizations.Options)) (*organizations.ListAccountsOutput, error) {
	if m.ListAccountsFunc != nil {
		return m.ListAccountsFunc(ctx, params, optFns...)
	}
	return nil, nil
}

// MockSTSAPI is a mock implementation of the STS API
type MockSTSAPI struct {
	mock.Mock
}

func (m *MockSTSAPI) AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sts.AssumeRoleOutput), args.Error(1)
}

func TestGetAccountInfo(t *testing.T) {
	// Create a mock client that returns test data
	mockClient := &mockOrganizationsClient{
		ListAccountsFunc: func(ctx context.Context, params *organizations.ListAccountsInput, optFns ...func(*organizations.Options)) (*organizations.ListAccountsOutput, error) {
			return &organizations.ListAccountsOutput{
				Accounts: []orgTypes.Account{
					{
						Id:     aws.String("123456789012"),
						Name:   aws.String("Test Account 1"),
						Email:  aws.String("test1@example.com"),
						Status: orgTypes.AccountStatusActive,
					},
					{
						Id:     aws.String("098765432109"),
						Name:   aws.String("Test Account 2"),
						Email:  aws.String("test2@example.com"),
						Status: orgTypes.AccountStatusActive,
					},
				},
			}, nil
		},
	}

	// Create a test client with the mock
	testClient := &Client{
		awsConfig: aws.Config{},
		orgClient: mockClient,
	}

	// Get account info
	accounts, err := testClient.GetAccountInfo(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, accounts)
	assert.Greater(t, len(accounts), 0)

	// Verify the first account's data
	if len(accounts) > 0 {
		assert.Equal(t, "123456789012", accounts[0].AccountId)
		assert.Equal(t, "Test Account 1", accounts[0].AccountName)
		assert.Equal(t, "test1@example.com", accounts[0].Email)
		assert.Equal(t, "ACTIVE", accounts[0].Status)
	}
}

func TestAssumeRole(t *testing.T) {
	// Create mock STS client
	mockSTS := new(MockSTSAPI)

	// Set up test data
	testCredentials := &stsTypes.Credentials{
		AccessKeyId:     aws.String("test-access-key"),
		SecretAccessKey: aws.String("test-secret-key"),
		SessionToken:    aws.String("test-session-token"),
		Expiration:      aws.Time(time.Now().Add(1 * time.Hour)),
	}

	// Set up expectations
	mockSTS.On("AssumeRole", mock.Anything, mock.AnythingOfType("*sts.AssumeRoleInput")).
		Return(&sts.AssumeRoleOutput{
			Credentials: testCredentials,
		}, nil)

	// Create test client with region and mock STS client
	testClient := &Client{
		awsConfig: aws.Config{
			Region: "us-west-2",
		},
		stsClient: mockSTS,
	}

	// Test assume role configuration
	assumeRoleConfig := &AssumeRoleConfig{
		RoleArn:         "arn:aws:iam::123456789012:role/TestRole",
		SessionName:     "test-session",
		DurationSeconds: 3600,
	}

	// Execute test
	err := testClient.AssumeRole(context.Background(), assumeRoleConfig)

	// Verify
	assert.NoError(t, err)
	mockSTS.AssertExpectations(t)

	// Verify credentials were updated
	creds, err := testClient.awsConfig.Credentials.Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "test-access-key", creds.AccessKeyID)
	assert.Equal(t, "test-secret-key", creds.SecretAccessKey)
	assert.Equal(t, "test-session-token", creds.SessionToken)
}
