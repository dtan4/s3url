package s3

import (
	"testing"

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
