package config

import (
	"testing"
)

func TestParseS3URL(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		url    string
		bucket string
		key    string
	}{
		{"https://s3-ap-northeast-1.amazonaws.com/bucket/key.txt", "bucket", "key.txt"},
		{"https://s3-ap-northeast-1.amazonaws.com/bucket/dir/key.txt", "bucket", "dir/key.txt"},
		{"https://bucket.s3.amazonaws.com/key.txt", "bucket", "key.txt"},
		{"https://bucket.s3.amazonaws.com/dir/key.txt", "bucket", "dir/key.txt"},
		{"https://bucket.s3-ap-northeast-1.amazonaws.com/key.txt", "bucket", "key.txt"},
		{"https://bucket.s3-ap-northeast-1.amazonaws.com/dir/key.txt", "bucket", "dir/key.txt"},
		{"s3://bucket/key.txt", "bucket", "key.txt"},
		{"s3://bucket/dir/key.txt", "bucket", "dir/key.txt"},
	}

	for _, tc := range testcases {
		c := &Config{}

		if err := c.ParseS3URL(tc.url); err != nil {
			t.Errorf("Error should not be raised. url: %s, error: %v", tc.url, err)
		}

		if c.Bucket != tc.bucket {
			t.Errorf("Bucket does not match. expected: %s, actual: %s", tc.bucket, c.Bucket)
		}

		if c.Key != tc.key {
			t.Errorf("Key does not match. expected: %s, actual: %s", tc.key, c.Key)
		}
	}
}

func TestParseURL_invalid(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		url    string
		errMsg string
	}{
		{
			url:    "foobarbaz",
			errMsg: "invalid hostname: url: \"foobarbaz\", hostname: \"\"",
		},
		{
			url:    "https://s3-ap-northeast-1.amazonaws.com/bucket",
			errMsg: "invalid path: url: \"https://s3-ap-northeast-1.amazonaws.com/bucket\", path: \"/bucket\"",
		},
	}

	for _, tc := range testcases {
		c := &Config{}

		err := c.ParseS3URL(tc.url)
		if err == nil {
			t.Error("Error should be raised.")
		}

		if err.Error() != tc.errMsg {
			t.Errorf("Error message does not match. expected: %s, actual: %s", tc.errMsg, err.Error())
		}
	}
}
