package s3

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dtan4/s3url/awsmock"
	"github.com/golang/mock/gomock"
)

func TestParseURL(t *testing.T) {
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
		bucket, key, err := ParseURL(tc.url)
		if err != nil {
			t.Errorf("Error should not be raised. url: %s, error: %v", tc.url, err)
		}

		if bucket != tc.bucket {
			t.Errorf("Bucket does not match. expected: %s, actual: %s", tc.bucket, bucket)
		}

		if key != tc.key {
			t.Errorf("Key does not match. expected: %s, actual: %s", tc.key, key)
		}
	}
}

func TestParseURL_invalid(t *testing.T) {
	testcases := []struct {
		url    string
		errMsg string
	}{
		{
			url:    "foobarbaz",
			errMsg: "Invalid URL, hostname is invalid. url: \"foobarbaz\", hostname: \"\"",
		},
		{
			url:    "https://s3-ap-northeast-1.amazonaws.com/bucket",
			errMsg: "Invalid URL, path is invalid. url: \"https://s3-ap-northeast-1.amazonaws.com/bucket\", path: \"/bucket\"",
		},
	}

	for _, tc := range testcases {
		_, _, err := ParseURL(tc.url)
		if err == nil {
			t.Error("Error should be raised.")
		}

		if err.Error() != tc.errMsg {
			t.Errorf("Error message does not match. expected: %s, actual: %s", tc.errMsg, err.Error())
		}
	}
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3mock := awsmock.NewMockS3API(ctrl)
	client := New(s3mock)

	if client.api != s3mock {
		t.Error("api does not match.")
	}
}

func TestGetPresignedURL(t *testing.T) {
	bucket := "bucket"
	key := "key"
	duration := int64(100)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := &url.URL{
		Scheme: "http",
		Host:   "bucket.s3-ap-northeast-1.amazonaws.com",
		Path:   "/key",
	}

	s3mock := awsmock.NewMockS3API(ctrl)
	s3mock.EXPECT().GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}).Return(&request.Request{
		HTTPRequest: &http.Request{
			URL: u,
		},
	}, &s3.GetObjectOutput{})
	client := &Client{
		api: s3mock,
	}

	signedURL, err := client.GetPresignedURL(bucket, key, duration)
	if err != nil {
		t.Fatalf("Error should not be raised. error:%s", err)
	}

	if signedURL != u.String() {
		t.Errorf("Invalid signed URL. signedURL: %s", signedURL)
	}
}

func TestUploadToS3(t *testing.T) {
	bucket := "bucket"
	key := "key"
	body := []byte("filebody")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3mock := awsmock.NewMockS3API(ctrl)
	s3mock.EXPECT().PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(body)),
	}).Return(&s3.PutObjectOutput{}, nil)
	client := &Client{
		api: s3mock,
	}

	if err := client.UploadToS3(bucket, key, body); err != nil {
		t.Fatalf("Error should not be raised. error: %s", err)
	}
}
