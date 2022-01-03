package s3bucket_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/s3bucket"

	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:   "S3Bucket",
			Golden: "s3bucket.yaml",
			Content: s3bucket.New(
				"myResource",
				"my-Bucket",
			),
		},
		{
			Name:   "S3Bucket with blocked public access",
			Golden: "s3bucket-public-block.yaml",
			Content: func() *s3bucket.S3Bucket {
				b := s3bucket.New(
					"myResource",
					"my-Bucket",
				)
				b.BlockAllPublicAccess = true

				return b
			}(),
		},
	}

	tstr.RunTests(t, testCases, nil)
}
