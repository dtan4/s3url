package s3

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dtan4/s3url/aws/mock"
	"github.com/golang/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3mock := mock.NewMockS3API(ctrl)
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

	s3mock := mock.NewMockS3API(ctrl)
	s3mock.EXPECT().GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}).Return(&request.Request{
		HTTPRequest: &http.Request{
			URL: u,
		},
		Operation: &request.Operation{},
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

	s3mock := mock.NewMockS3API(ctrl)
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
