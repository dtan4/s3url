package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getPresignedURL(svc *s3.S3, bucket, key string, duration int64) (string, error) {
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	signedURL, err := req.Presign(time.Duration(duration) * time.Minute)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func parseURL(s3URL string) (string, string, error) {
	var bucket, key string

	u, err := url.Parse(s3URL)
	if err != nil {
		return "", "", fmt.Errorf("Invalid URL: %s.\n", s3URL)
	}

	if u.Scheme == "s3" { // s3://bucket/key
		bucket = u.Host
		key = strings.Replace(u.Path, "/", "", 1)
	} else { // https://s3-ap-northeast-1.amazonaws.com/bucket/key
		ss := strings.SplitN(u.Path, "/", 3)
		bucket = ss[1]
		key = ss[2]
	}

	return bucket, key, nil
}

func uploadToS3(svc *s3.S3, path, bucket, key string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   fp,
	})
	if err != nil {
		return err
	}

	return nil
}
