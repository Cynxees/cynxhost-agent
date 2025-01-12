package dependencies

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type AWSClient struct {
	Config    *aws.Config
	EC2Client *ec2.Client
}

func NewAWSClient(accessKeyId string, secret string) *AWSClient {

	config := newAWSConfig(accessKeyId, secret)

	return &AWSClient{
		Config:    config,
		EC2Client: newEC2Client(*config),
	}
}

func newEC2Client(config aws.Config) *ec2.Client {
	return ec2.NewFromConfig(config)
}

func newAWSConfig(accessKeyId string, secret string) *aws.Config {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKeyId,
				secret,
				"", // Optional session token (use if MFA is enabled, otherwise leave it empty)
			),
		),
		config.WithRegion("ap-southeast-1"),
	)

	if err != nil {
		panic(err)
	}

	return &cfg
}
