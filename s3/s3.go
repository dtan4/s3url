package s3

import (
	"bytes"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// Client represents the wrapper of S3 API Client
type Client struct {
	api s3iface.S3API
}

// NewClient creates new Client
func NewClient(api s3iface.S3API) *Client {
	return &Client{
		api: api,
	}
}

// GetPresignedURL returns S3 object pre-signed URL
func (c *Client) GetPresignedURL(bucket, key string, duration int64) (string, error) {
	req, _ := c.api.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	signedURL, err := req.Presign(time.Duration(duration) * time.Minute)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

// UploadToS3 uploads local file to the specified S3 location
func (c *Client) UploadToS3(bucket, key string, body []byte) error {
	_, err := c.api.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return err
	}

	return nil
}
