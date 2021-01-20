package s3

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type fakeS3Client struct {
	putObject func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func (c *fakeS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return c.putObject(ctx, params, optFns...)
}

type fakeS3PresignClient struct {
	presignGetObject func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

func (c *fakeS3PresignClient) PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	return c.presignGetObject(ctx, params, optFns...)
}

func TestNew(t *testing.T) {
	t.Parallel()

	s3Client := &fakeS3Client{}
	s3PresignClient := &fakeS3PresignClient{}
	client := New(s3Client, s3PresignClient)

	if client.s3Client != s3Client {
		t.Error("s3Client does not match")
	}

	if client.s3PresignClient != s3PresignClient {
		t.Error("s3Client does not match")
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
		s3PresignClient: &fakeS3PresignClient{
			presignGetObject: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
				u := &url.URL{
					Scheme: "http",
					Host:   fmt.Sprintf("%s.s3-ap-northeast-1.amazonaws.com", aws.ToString(params.Bucket)),
					Path:   fmt.Sprintf("/%s", aws.ToString(params.Key)),
				}

				return &v4.PresignedHTTPRequest{
					URL: u.String(),
				}, nil
			},
		},
	}

	got, err := client.GetPresignedURL(context.Background(), bucket, key, duration)
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
		s3Client: &fakeS3Client{
			putObject: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
		},
	}

	if err := client.UploadToS3(context.Background(), bucket, key, f); err != nil {
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
