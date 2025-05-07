package awshelper

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type ifConfigService interface {
	LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (cfg aws.Config, err error)
}

type ConfigService struct {
	Service ifConfigService
}

type configLoader struct{}

func (c *configLoader) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (cfg aws.Config, err error) {
	return config.LoadDefaultConfig(ctx, optFns...)
}

func NewConfigService() ConfigService {
	return ConfigService{Service: &configLoader{}}
}

func (configService *ConfigService) FindAWSCredential(profile string) (aws.Config, error) {
	awsProfile := profile

	if os.Getenv("AWS_SESSION_TOKEN") == "" {
		if awsProfile == "default" && os.Getenv("AWS_DEFAULT_PROFILE") != "" {
			awsProfile = os.Getenv("AWS_DEFAULT_PROFILE")
		}
		awsCfg, err := configService.Service.LoadDefaultConfig(
			context.TODO(),
			config.WithSharedConfigProfile(awsProfile),
		)
		return awsCfg, err
	} else {
		awsCfg, err := configService.Service.LoadDefaultConfig(context.TODO())
		return awsCfg, err
	}
}
