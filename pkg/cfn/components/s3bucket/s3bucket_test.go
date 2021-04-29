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
	}

	tstr.RunTests(t, testCases, nil)
}
