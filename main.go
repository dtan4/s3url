package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3api "github.com/aws/aws-sdk-go/service/s3"
	"github.com/dtan4/s3url/s3"
	flag "github.com/spf13/pflag"
)

const (
	defaultDuration = 5
)

func main() {
	var (
		bucket   string
		duration int64
		key      string
		profile  string
		upload   string
		version  bool
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

	f.StringVarP(&bucket, "bucket", "b", "", "Bucket name")
	f.Int64VarP(&duration, "duration", "d", defaultDuration, "Valid duration in minutes")
	f.StringVarP(&key, "key", "k", "", "Object key")
	f.StringVar(&profile, "profile", "", "AWS profile name")
	f.StringVar(&upload, "upload", "", "File to upload")
	f.BoolVarP(&version, "version", "v", false, "Print version")

	f.Parse(os.Args[1:])

	if version {
		printVersion()
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

	if s3URL == "" && (bucket == "" || key == "") {
		f.Usage()
		os.Exit(1)
	}

	if s3URL != "" {
		bucket, key, err = s3.ParseURL(s3URL)
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

	api := s3api.New(sess, &aws.Config{})
	s3Client := s3.NewClient(api)

	if upload != "" {
		path, err := filepath.Abs(upload)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		body, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := s3Client.UploadToS3(bucket, key, body); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stderr, "uploaded: "+path)
	}

	signedURL, err := s3Client.GetPresignedURL(bucket, key, duration)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(signedURL)
}
