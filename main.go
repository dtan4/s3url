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

	"github.com/dtan4/s3url/config"
)

const (
	defaultDuration = 5
)

func run(args []string) error {
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

	c := config.Config{}

	f.StringVarP(&c.Bucket, "bucket", "b", "", "Bucket name")
	f.Int64VarP(&c.Duration, "duration", "d", defaultDuration, "Valid duration in minutes")
	f.StringVarP(&c.Key, "key", "k", "", "Object key")
	f.StringVar(&c.Profile, "profile", "", "AWS profile name")
	f.StringVar(&c.Upload, "upload", "", "File to upload")
	f.BoolVarP(&c.Version, "version", "v", false, "Print version")

	f.Parse(args[1:])

	if c.Version {
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

	if c.Profile != "" {
		sess, err = session.NewSessionWithOptions(session.Options{
			Profile: c.Profile,
		})
		if err != nil {
			return err
		}
	} else {
		sess = session.New()
	}

	if s3URL == "" && (c.Bucket == "" || c.Key == "") {
		f.Usage()
		return fmt.Errorf("insufficient arguments")
	}

	if s3URL != "" {
		c.Bucket, c.Key, err = s3.ParseURL(s3URL)
		if err != nil {
			return err
		}
	}

	if c.Bucket == "" {
		return fmt.Errorf("Bucket name is required.")
	}

	if c.Key == "" {
		return fmt.Errorf("Object key is required.")
	}

	api := s3api.New(sess, &aws.Config{})
	s3Client := s3.New(api)

	if c.Upload != "" {
		path, err := filepath.Abs(c.Upload)
		if err != nil {
			return err
		}

		body, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if err := s3Client.UploadToS3(c.Bucket, c.Key, body); err != nil {
			return err
		}

		fmt.Fprintln(os.Stderr, "uploaded: "+path)
	}

	signedURL, err := s3Client.GetPresignedURL(c.Bucket, c.Key, c.Duration)
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
