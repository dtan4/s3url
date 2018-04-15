package config

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	// DefaultDuration represents default valid duration in minutes
	DefaultDuration = 5
)

var (
	virtualHostRegexp = regexp.MustCompile(`^s3-[a-z0-9-]+\.amazonaws\.com$`)
)

// Config represents s3url configurations
type Config struct {
	Bucket   string
	Duration int64
	Key      string
	Profile  string
	Upload   string
	Version  bool
}

// ParseS3URL extracts bucket and key from S3 URL
func (c *Config) ParseS3URL(s3URL string) error {
	u, err := url.Parse(s3URL)
	if err != nil {
		return errors.Wrapf(err, "invalid URL: %s", s3URL)
	}

	if u.Scheme == "s3" { // s3://bucket/key
		c.Bucket = u.Host
		c.Key = strings.Replace(u.Path, "/", "", 1)
	} else {
		if virtualHostRegexp.MatchString(u.Host) { // https://s3-ap-northeast-1.amazonaws.com/bucket/key
			ss := strings.SplitN(u.Path, "/", 3)
			if len(ss) < 3 {
				return errors.Errorf("invalid path: url: %q, path: %q", s3URL, u.Path)
			}

			c.Bucket = ss[1]
			c.Key = ss[2]
		} else { // https://bucket.s3-ap-northeast-1.amazonaws.com/key
			ss := strings.Split(u.Host, ".")
			if len(ss) < 4 {
				return errors.Errorf("invalid hostname: url: %q, hostname: %q", s3URL, u.Host)
			}

			c.Bucket = strings.Join(ss[0:len(ss)-3], ".")
			c.Key = u.Path[1:]
		}
	}

	return nil
}
