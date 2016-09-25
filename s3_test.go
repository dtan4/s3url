package main

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	testdata := []struct {
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

	for _, tt := range testdata {
		bucket, key, err := parseURL(tt.url)
		if err != nil {
			t.Errorf("Error should not be raised. url: %s, error: %v", tt.url, err)
		}

		if bucket != tt.bucket {
			t.Errorf("Bucket does not matched. expect: %s, actual: %s", tt.bucket, bucket)
		}

		if key != tt.key {
			t.Errorf("Key does not matched. expect: %s, actual: %s", tt.key, key)
		}
	}
}
