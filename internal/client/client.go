package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	stsTypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
)

// Client represents the AWS client with assume role support
type Client struct {
	awsConfig aws.Config
}

// AssumeRoleConfig represents the configuration for assuming a role
type AssumeRoleConfig struct {
	RoleArn           string
	SessionName       string
	ExternalId        string
	DurationSeconds   int32
	PolicyArns        []string
	Policy            string
	Tags              map[string]string
	TransitiveTagKeys []string
}

// NewClient creates a new AWS client with the given credentials
func NewClient(accessKey, secretKey, region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &Client{
		awsConfig: cfg,
	}, nil
}

// AssumeRole assumes the specified IAM role and returns new credentials
func (c *Client) AssumeRole(ctx context.Context, assumeRoleConfig *AssumeRoleConfig) error {
	stsClient := sts.NewFromConfig(c.awsConfig)

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(assumeRoleConfig.RoleArn),
		RoleSessionName: aws.String(assumeRoleConfig.SessionName),
	}

	if assumeRoleConfig.ExternalId != "" {
		input.ExternalId = aws.String(assumeRoleConfig.ExternalId)
	}

	if assumeRoleConfig.DurationSeconds > 0 {
		input.DurationSeconds = aws.Int32(assumeRoleConfig.DurationSeconds)
	}

	if len(assumeRoleConfig.PolicyArns) > 0 {
		policyArns := make([]stsTypes.PolicyDescriptorType, len(assumeRoleConfig.PolicyArns))
		for i, arn := range assumeRoleConfig.PolicyArns {
			policyArns[i] = stsTypes.PolicyDescriptorType{
				Arn: aws.String(arn),
			}
		}
		input.PolicyArns = policyArns
	}

	if assumeRoleConfig.Policy != "" {
		input.Policy = aws.String(assumeRoleConfig.Policy)
	}

	if len(assumeRoleConfig.Tags) > 0 {
		tags := make([]stsTypes.Tag, 0, len(assumeRoleConfig.Tags))
		for k, v := range assumeRoleConfig.Tags {
			tags = append(tags, stsTypes.Tag{
				Key:   aws.String(k),
				Value: aws.String(v),
			})
		}
		input.Tags = tags
	}

	if len(assumeRoleConfig.TransitiveTagKeys) > 0 {
		input.TransitiveTagKeys = assumeRoleConfig.TransitiveTagKeys
	}

	result, err := stsClient.AssumeRole(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to assume role: %w", err)
	}

	// Update the client's credentials with the assumed role credentials
	c.awsConfig.Credentials = credentials.NewStaticCredentialsProvider(
		aws.ToString(result.Credentials.AccessKeyId),
		aws.ToString(result.Credentials.SecretAccessKey),
		aws.ToString(result.Credentials.SessionToken),
	)

	return nil
}

// GetAccountInfo retrieves information about AWS accounts from AWS Organizations
func (c *Client) GetAccountInfo(ctx context.Context) ([]AccountInfo, error) {
	orgClient := organizations.NewFromConfig(c.awsConfig)

	var accounts []AccountInfo
	var nextToken *string

	for {
		input := &organizations.ListAccountsInput{
			NextToken: nextToken,
		}

		result, err := orgClient.ListAccounts(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list accounts: %w", err)
		}

		for _, account := range result.Accounts {
			accounts = append(accounts, AccountInfo{
				AccountId:   aws.ToString(account.Id),
				AccountName: aws.ToString(account.Name),
				Email:       aws.ToString(account.Email),
				Status:      string(account.Status),
			})
		}

		if result.NextToken == nil {
			break
		}
		nextToken = result.NextToken
	}

	return accounts, nil
}

// AccountInfo represents information about an AWS account
type AccountInfo struct {
	AccountId   string
	AccountName string
	Email       string
	Status      string
}

// SetSessionToken sets the AWS session token for temporary credentials
func (c *Client) SetSessionToken(token string) error {
	// Get current credentials
	creds, err := c.awsConfig.Credentials.Retrieve(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %w", err)
	}

	// Create new credentials with the session token
	c.awsConfig.Credentials = credentials.NewStaticCredentialsProvider(
		creds.AccessKeyID,
		creds.SecretAccessKey,
		token,
	)
	return nil
}
