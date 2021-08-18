package s3

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type s3PresignClient interface {
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// Client represents the wrapper of S3 API Client
type Client struct {
	s3Client        s3Client
	s3PresignClient s3PresignClient
}

// New creates new Client
func New(s3Client s3Client, s3PresignClient s3PresignClient) *Client {
	return &Client{
		s3Client:        s3Client,
		s3PresignClient: s3PresignClient,
	}
}

// GetPresignedURL returns S3 object pre-signed URL
func (c *Client) GetPresignedURL(ctx context.Context, bucket, key string, duration int64) (string, error) {
	req, err := c.s3PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(time.Duration(duration)*time.Minute))
	if err != nil {
		return "", fmt.Errorf("cannot generate signed URL: %w", err)
	}

	return req.URL, nil
}

// UploadToS3 uploads local file to the specified S3 location
func (c *Client) UploadToS3(ctx context.Context, bucket, key string, reader io.ReadSeeker) error {
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("cannot upload file to S3: %w", err)
	}

	return nil
}
