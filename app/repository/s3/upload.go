package s3_repository

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (r *s3Repo) UploadFilePublic(ctx context.Context, objectName string, body io.Reader, mimeType string) (string, error) {
	// upload to bucket s3
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(objectName),
		Body:        body,
		ContentType: aws.String(mimeType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s/%s", os.Getenv("S3_ENDPOINT"), r.bucketName, objectName)
	return url, nil
}
