package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
		help     bool
		key      string
		profile  string
		upload   string
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
	f.BoolVar(&help, "help", false, "Show usage")
	f.BoolVar(&help, "h", false, "Show usage")
	f.StringVar(&key, "key", "", "Object key")
	f.StringVar(&key, "k", "", "Object key")
	f.StringVar(&profile, "profile", "", "AWS profile name")
	f.StringVar(&upload, "upload", "", "File to upload")

	f.Parse(os.Args[1:])

	if help {
		f.Usage()
		os.Exit(0)
	}

	var s3URL string

	for 0 < f.NArg() {
		s3URL = f.Args()[0]
		f.Parse(f.Args()[1:])
	}

	var sess *session.Session
	var err error

	if profile != "" {
		sess, err = session.NewSessionWithOptions(session.Options{
			Profile: profile,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		sess = session.New()
	}

	if s3URL != "" {
		bucket, key, err = parseURL(s3URL)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if bucket == "" {
		fmt.Fprintln(os.Stderr, "Bucket name is required.")
		os.Exit(1)
	}

	if key == "" {
		fmt.Fprintln(os.Stderr, "Object key is required.")
		os.Exit(1)
	}

	svc := s3.New(sess, &aws.Config{})

	if upload != "" {
		path, err := filepath.Abs(upload)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := uploadToS3(svc, path, bucket, key); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stderr, "uploaded: "+path)
	}

	signedURL, err := getPresignedURL(svc, bucket, key, duration)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(signedURL)
}
