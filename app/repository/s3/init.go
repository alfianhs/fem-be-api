package s3_repository

import (
	"app/domain"
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

type s3Repo struct {
	bucketName string
	client     *s3.Client
}

func NewS3Repository(contextTimeout time.Duration) domain.S3Repo {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     os.Getenv("S3_ACCESS_KEY"),
				SecretAccessKey: os.Getenv("S3_SECRET_KEY"),
				SessionToken:    "",
			},
		}),
		config.WithRegion(os.Getenv("S3_REGION")),
		config.WithBaseEndpoint(os.Getenv("S3_ENDPOINT")),
	)
	if err != nil {
		logrus.Error("Load S3 Config Error : ", err)
		return nil
	}

	client := s3.NewFromConfig(cfg)
	return &s3Repo{
		client:     client,
		bucketName: os.Getenv("S3_BUCKET_NAME"),
	}
}
