package s3

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dtan4/s3url/aws/mock"
	"github.com/golang/mock/gomock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3mock := mock.NewMockS3API(ctrl)
	client := New(s3mock)

	if client.api != s3mock {
		t.Error("api does not match.")
	}
}

func TestGetPresignedURL(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	bucket := "bucket"
	key := "key"
	testfile := filepath.Join("..", "..", "_testdata", "test.txt")

	f, err := os.Open(testfile)
	if err != nil {
		t.Fatalf("cannot open testdata %q: %s", testfile, err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3mock := mock.NewMockS3API(ctrl)
	// TODO: hard to write io.ReadSeeker expectation
	s3mock.EXPECT().PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   f,
	}).Return(&s3.PutObjectOutput{}, nil)
	client := &Client{
		api: s3mock,
	}

	if err := client.UploadToS3(bucket, key, f); err != nil {
		t.Fatalf("Error should not be raised. error: %s", err)
	}
}

func BenchmarkReadFileEntirely(b *testing.B) {
	testfile := filepath.Join("..", "..", "_testdata", "test.txt")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		body, _ := ioutil.ReadFile(testfile)
		r := bytes.NewReader(body)
		_, _ = ioutil.ReadAll(r)
	}
}

func BenchmarkReadFileStream(b *testing.B) {
	testfile := filepath.Join("..", "..", "_testdata", "test.txt")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f, _ := os.Open(testfile)
		_, _ = ioutil.ReadAll(f)
		f.Close()
	}
}
