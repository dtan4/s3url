package s3

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dtan4/s3url/awsmock"
	"github.com/golang/mock/gomock"
)

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3mock := awsmock.NewMockS3API(ctrl)
	client := NewClient(s3mock)

	if client.api != s3mock {
		t.Error("api does not match.")
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
		t.Errorf("Error should not be raised. error: %s", err)
	}
}
