package s3

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type mockS3API struct {
	s3iface.S3API
}

func (m *mockS3API) GetObjectRequest(input *s3.GetObjectInput) (*request.Request, *s3.GetObjectOutput) {
	return &request.Request{
		HTTPRequest: &http.Request{
			URL: &url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("%s.s3-ap-northeast-1.amazonaws.com", aws.StringValue(input.Bucket)),
				Path:   fmt.Sprintf("/%s", aws.StringValue(input.Key)),
			},
		},
		Operation: &request.Operation{},
	}, &s3.GetObjectOutput{}
}

func (m *mockS3API) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return &s3.PutObjectOutput{}, nil
}

func TestNew(t *testing.T) {
	t.Parallel()

	api := &mockS3API{}
	client := New(api)

	if client.api != api {
		t.Error("api does not match.")
	}
}

func TestGetPresignedURL(t *testing.T) {
	t.Parallel()

	bucket := "bucket"
	key := "key"
	duration := int64(100)

	u := url.URL{
		Scheme: "http",
		Host:   "bucket.s3-ap-northeast-1.amazonaws.com",
		Path:   "/key",
	}
	want := u.String()

	client := &Client{
		api: &mockS3API{},
	}

	got, err := client.GetPresignedURL(bucket, key, duration)
	if err != nil {
		t.Fatalf("Error should not be raised. error:%s", err)
	}

	if got != want {
		t.Errorf("Invalid signed URL. want: %s, got: %s", want, got)
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

	client := &Client{
		api: &mockS3API{},
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
