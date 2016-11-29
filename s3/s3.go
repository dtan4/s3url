package s3

import (
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
