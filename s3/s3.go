package s3

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var (
	virtualHostRegexp = regexp.MustCompile(`^s3-[a-z0-9-]+\.amazonaws\.com$`)
)

// Client represents the wrapper of S3 API Client
type Client struct {
	api s3iface.S3API
}

// ParseURL parses S3 object URL
func ParseURL(s3URL string) (string, string, error) {
	var bucket, key string

	u, err := url.Parse(s3URL)
	if err != nil {
		return "", "", fmt.Errorf("Invalid URL: %s", s3URL)
	}

	if u.Scheme == "s3" { // s3://bucket/key
		bucket = u.Host
		key = strings.Replace(u.Path, "/", "", 1)
	} else {
		if virtualHostRegexp.MatchString(u.Host) { // https://s3-ap-northeast-1.amazonaws.com/bucket/key
			ss := strings.SplitN(u.Path, "/", 3)
			if len(ss) < 3 {
				return "", "", fmt.Errorf("Invalid path: %s", u.Path)
			}

			bucket = ss[1]
			key = ss[2]
		} else { // https://bucket.s3-ap-northeast-1.amazonaws.com/key
			ss := strings.Split(u.Host, ".")
			if len(ss) < 4 {
				return "", "", fmt.Errorf("Invalid hostname: %s", u.Host)
			}

			bucket = strings.Join(ss[0:len(ss)-3], ".")
			key = u.Path[1:]
		}
	}

	return bucket, key, nil
}

// New creates new Client
func New(api s3iface.S3API) *Client {
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
