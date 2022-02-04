package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

//nolint:funlen
func TestCustomValidateFieldCanNotContainString(t *testing.T) {
	testCases := []struct {
		name                string
		fieldValue          string
		canNotContainString string
		expectErr           string
	}{
		{
			name:                "field doesn't contain a invalid string",
			fieldValue:          "administrator",
			canNotContainString: "foobar",
		},
		{
			name:                "field contains a invalid string",
			fieldValue:          "double--hyphen",
			canNotContainString: "--",
			expectErr:           "Field can not contain a double hyphen",
		},
		{
			name:                "field consists of a single invalid character",
			fieldValue:          "$",
			canNotContainString: "$",
			expectErr:           "Field can not contain the character '$'",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			validateFunc := v1alpha1.ValidateFieldCanNotContainString(tc.canNotContainString, tc.expectErr)
			err := validateFunc(tc.fieldValue)
			if tc.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.expectErr, err.Error())
			}
		})
	}
}
