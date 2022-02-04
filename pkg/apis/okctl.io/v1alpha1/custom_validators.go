package v1alpha1

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mishudark/errors"
)

// ValidateFieldCanNotContainString Check if a validation.Field contains a illegal string
// Usage:
// validation.Field(&struct.Name,
// 	validation.By(ValidateFieldCanNotContainString("--", "field can not have two consecutive hyphens")),
// ),
func ValidateFieldCanNotContainString(str string, errorString string) validation.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
		res := strings.Contains(s, str)

		if res {
			return errors.New(errorString)
		}

		return nil
	}
}
