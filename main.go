package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	defaultDuration = 5
)

func main() {
	var (
		bucket   string
		duration int64
		key      string
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage of %s:
   %s [OPTIONS]

Options:
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&bucket, "bucket", "", "Bucket name")
	flag.StringVar(&bucket, "b", "", "Bucket name")
	flag.Int64Var(&duration, "duration", defaultDuration, "Valid duration in minutes")
	flag.Int64Var(&duration, "d", defaultDuration, "Valid duration in minutes")
	flag.StringVar(&key, "key", "", "Object key")
	flag.StringVar(&key, "k", "", "Object key")

	flag.Parse()

	if bucket == "" {
		fmt.Fprintln(os.Stderr, "Bucket name is required.")
		os.Exit(1)
	}

	if key == "" {
		fmt.Fprintln(os.Stderr, "Object key is required.")
		os.Exit(1)
	}

	svc := s3.New(session.New(), &aws.Config{})
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(time.Duration(duration) * time.Minute)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(url)
}
