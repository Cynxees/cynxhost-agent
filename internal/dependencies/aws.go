package dependencies

import (
	"context"
	"cynxhostagent/internal/constant"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

type AWSClient struct {
	Config    *aws.Config
	EC2Client *ec2.Client
	ECRClient *ecr.Client
}

func NewAWSClient(accessKeyId string, secret string) *AWSClient {

	config := newAWSConfig(accessKeyId, secret)

	return &AWSClient{
		Config:    config,
		EC2Client: newEC2Client(*config),
		ECRClient: newECRClient(*config),
	}
}

func newEC2Client(config aws.Config) *ec2.Client {
	return ec2.NewFromConfig(config)
}

func newECRClient(config aws.Config) *ecr.Client {
	return ecr.NewFromConfig(config)
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

func (a *AWSClient) GetRepository(ctx context.Context, repositoryName string) (*ecr.DescribeRepositoriesOutput, error) {
	input := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{repositoryName},
	}

	out, err := a.ECRClient.DescribeRepositories(ctx, input)
	if err != nil {

		if strings.Contains(err.Error(), "RepositoryNotFoundException") {
			return out, errors.New(constant.NotFound)
		}

		return out, err
	}

	log.Println("Repository exists")
	return out, nil
}

func (a *AWSClient) CreateRepository(ctx context.Context, repositoryName string) (*ecr.CreateRepositoryOutput, error) {
	input := &ecr.CreateRepositoryInput{
		RepositoryName: aws.String(repositoryName),
	}

	out, err := a.ECRClient.CreateRepository(ctx, input)
	if err != nil {
		return out, err
	}

	log.Println("Repository created")
	return out, nil
}
