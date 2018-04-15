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

func run(args []string) error {
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
`, args[0], args[0], args[0], args[0])
		f.PrintDefaults()
	}

	f.StringVarP(&bucket, "bucket", "b", "", "Bucket name")
	f.Int64VarP(&duration, "duration", "d", defaultDuration, "Valid duration in minutes")
	f.StringVarP(&key, "key", "k", "", "Object key")
	f.StringVar(&profile, "profile", "", "AWS profile name")
	f.StringVar(&upload, "upload", "", "File to upload")
	f.BoolVarP(&version, "version", "v", false, "Print version")

	f.Parse(args[1:])

	if version {
		printVersion()
		return nil
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
			return err
		}
	} else {
		sess = session.New()
	}

	if s3URL == "" && (bucket == "" || key == "") {
		f.Usage()
		return fmt.Errorf("insufficient arguments")
	}

	if s3URL != "" {
		bucket, key, err = s3.ParseURL(s3URL)
		if err != nil {
			return err
		}
	}

	if bucket == "" {
		return fmt.Errorf("Bucket name is required.")
	}

	if key == "" {
		return fmt.Errorf("Object key is required.")
	}

	api := s3api.New(sess, &aws.Config{})
	s3Client := s3.New(api)

	if upload != "" {
		path, err := filepath.Abs(upload)
		if err != nil {
			return err
		}

		body, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if err := s3Client.UploadToS3(bucket, key, body); err != nil {
			return err
		}

		fmt.Fprintln(os.Stderr, "uploaded: "+path)
	}

	signedURL, err := s3Client.GetPresignedURL(bucket, key, duration)
	if err != nil {
		return err
	}

	fmt.Println(signedURL)

	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
