package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3api "github.com/aws/aws-sdk-go/service/s3"

	"github.com/dtan4/s3url/aws/s3"
)

var (
	// S3 represents S3 API client
	S3 *s3.Client
)

// Initialize creates S3 API client objects
func Initialize(profile string) error {
	var (
		sess *session.Session
		err  error
	)

	if profile != "" {
		sess, err = session.NewSessionWithOptions(session.Options{
			Profile: profile,
		})
		if err != nil {
			return err
		}
	} else {
		sess = session.New()
	}

	api := s3api.New(sess, &aws.Config{})
	S3 = s3.New(api)

	return nil
}
