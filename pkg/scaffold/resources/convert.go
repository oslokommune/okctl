// Package scaffold knows how to scaffold okctl applications
package resources

import (
	"bytes"
	"fmt"

	kaex "github.com/oslokommune/kaex/pkg/api"
	v1 "k8s.io/api/core/v1"
)

func VolumesAsBytes(volumes []*v1.PersistentVolumeClaim) ([]byte, error) {
	var writer bytes.Buffer

	for index := range volumes {
		err := kaex.WriteCleanResource(&writer, volumes[index])
		if err != nil {
			return nil, fmt.Errorf("writing volume to buffer: %w", err)
		}
	}

	return writer.Bytes(), nil
}

func ResourceAsBytes(data interface{}) ([]byte, error) {
	var writer bytes.Buffer

	err := kaex.WriteCleanResource(&writer, data)
	if err != nil {
		return nil, fmt.Errorf("writing resource definition to buffer: %w", err)
	}

	return writer.Bytes(), nil
}
