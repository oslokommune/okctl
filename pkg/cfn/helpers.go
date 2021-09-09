package cfn

import (
	"fmt"

	cfPkg "github.com/aws/aws-sdk-go/service/cloudformation"
)

// GetOutput retrieves a certain output from a cfn describe stack response
func GetOutput(outputs []*cfPkg.Output, name string) (string, error) {
	for _, output := range outputs {
		if *output.OutputKey == name {
			return *output.OutputValue, nil
		}
	}

	return "", fmt.Errorf("output with key %s not found", name)
}
