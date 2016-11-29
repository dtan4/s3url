package main

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	testcases := []struct {
		url    string
		bucket string
		key    string
	}{
		{"https://s3-region.amazonaws.com/bucket/key.txt", "bucket", "key.txt"},
		{"https://s3-region.amazonaws.com/bucket/dir/key.txt", "bucket", "dir/key.txt"},
		{"https://bucket.s3.amazonaws.com/key.txt", "bucket", "key.txt"},
		{"https://bucket.s3.amazonaws.com/dir/key.txt", "bucket", "dir/key.txt"},
		{"https://bucket.s3-region.amazonaws.com/key.txt", "bucket", "key.txt"},
		{"https://bucket.s3-region.amazonaws.com/dir/key.txt", "bucket", "dir/key.txt"},
		{"s3://bucket/key.txt", "bucket", "key.txt"},
		{"s3://bucket/dir/key.txt", "bucket", "dir/key.txt"},
	}

	for _, tc := range testcases {
		bucket, key, err := parseURL(tc.url)
		if err != nil {
			t.Errorf("Error should not be raised. url: %s, error: %v", tc.url, err)
		}

		if bucket != tc.bucket {
			t.Errorf("Bucket does not matched. expect: %s, actual: %s", tc.bucket, bucket)
		}

		if key != tc.key {
			t.Errorf("Key does not matched. expect: %s, actual: %s", tc.key, key)
		}
	}
}
