package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

func (cfg *config) LoadConfig() {

	cfg.port = os.Getenv("PORT")
	if cfg.port == "" {
		os.Setenv("PORT", "50051")
		cfg.port = "50051"
		logger.Info("PORT env variable is not set, defaulting to 50051")
	}

	cfg.sender_email = os.Getenv("SENDER_EMAIL")
	if cfg.sender_email == "" {
		logger.Error("SENDER_MAIL env variable is not set")
		panic("SENDER_MAIL env variable is not set")
	}

	cfg.env = os.Getenv("ENV")
	if cfg.env == "" {
		os.Setenv("ENV", "dev")
		cfg.env = "dev"
		logger.Info("ENV env variable is not set, defaulting to dev")
	}

}

func (cfg *config) ConfigLocalAws() (aws.Config, error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {

		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:4566",
			SigningRegion: "us-east-1",
		}, nil

	})
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		logger.Error("Cannot load the AWS configs: %s", err)
		return awsCfg, err
	}
	return awsCfg, nil

}

func (c *config) ConfigAws() {
	var cfg aws.Config
	var err error

	if c.env == "dev" {
		cfg, err = c.ConfigLocalAws()
		if err != nil {
			log.Fatal(err.Error())

		}
	} else {
		cfg, err = awsconfig.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Fatal(err.Error())

		}
	}

	c.awsConfig = cfg

}
