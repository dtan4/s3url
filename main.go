package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	defaultDuration = 5
)

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

func main() {
	var (
		bucket   string
		duration int64
		key      string
	)

	f := flag.NewFlagSet("s3url", flag.ExitOnError)

	f.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage of %s:
   %s https://s3-region.amazonaws.com/BUCKET/KEY [-d DURATION]
   %s s3://BUCKET/KEY [-d DURATION]
   %s -b BUCKET -k KEY [-d DURATION]

Options:
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
		f.PrintDefaults()
	}

	f.StringVar(&bucket, "bucket", "", "Bucket name")
	f.StringVar(&bucket, "b", "", "Bucket name")
	f.Int64Var(&duration, "duration", defaultDuration, "Valid duration in minutes")
	f.Int64Var(&duration, "d", defaultDuration, "Valid duration in minutes")
	f.StringVar(&key, "key", "", "Object key")
	f.StringVar(&key, "k", "", "Object key")

	f.Parse(os.Args[1:])

	var s3URL string

	for 0 < f.NArg() {
		s3URL = f.Args()[0]
		f.Parse(f.Args()[1:])
	}

	if s3URL != "" {
		b, k, err := parseURL(s3URL)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		bucket, key = b, k
	}

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

	signedURL, err := req.Presign(time.Duration(duration) * time.Minute)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(signedURL)
}
