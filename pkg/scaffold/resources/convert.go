// Package resources knows how to produce and handle Kubernetes resource
package resources

import (
	"bytes"
	"fmt"

	kaex "github.com/oslokommune/kaex/pkg/api"
	v1 "k8s.io/api/core/v1"
)

// VolumesAsBytes knows how to convert a Kubernetes PersistentVolumeClaim to a byte array
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

// ResourceAsBytes knows how to convert a Kubernetes resource to a byte array
func ResourceAsBytes(data interface{}) ([]byte, error) {
	var writer bytes.Buffer

	err := kaex.WriteCleanResource(&writer, data)
	if err != nil {
		return nil, fmt.Errorf("writing resource definition to buffer: %w", err)
	}

	return writer.Bytes(), nil
}
