package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	s3svc "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/dtan4/s3url/aws/s3"
)

var (
	// S3 represents S3 API client
	S3 *s3.Client
)

// Initialize creates S3 API client objects
func Initialize(ctx context.Context, profile string) error {
	var (
		cfg aws.Config
		err error
	)

	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
		if err != nil {
			return fmt.Errorf("cannot load config using profile %q: %w", profile, err)
		}
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
		if err != nil {
			return fmt.Errorf("cannot load default config: %w", err)
		}
	}

	s3Client := s3svc.NewFromConfig(cfg)
	s3PresignClient := s3svc.NewPresignClient(s3Client)

	S3 = s3.New(s3Client, s3PresignClient)

	return nil
}
