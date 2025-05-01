package s3_repository

import (
	s3_model "app/domain/model/s3"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (r *s3Repo) UploadFilePublic(ctx context.Context, objectName string, body io.Reader, mimeType string) (*s3_model.UploadResponse, error) {
	// upload to bucket s3
	publicPrefix := "public/"
	fullObjectName := publicPrefix + objectName

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(fullObjectName),
		Body:        body,
		ContentType: aws.String(mimeType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/%s", os.Getenv("S3_ENDPOINT"), r.bucketName, fullObjectName)
	s3Response := &s3_model.UploadResponse{
		Key:       fullObjectName,
		URL:       url,
		ExpiredAt: nil,
	}
	return s3Response, nil
}
