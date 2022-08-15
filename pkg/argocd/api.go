package argocd

import (
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/scaffold/resources"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"
)

//nolint:lll
const argoCDUnmarshallingError = "cannot unmarshal string into Go value of type struct { APIVersion string \"json:\\\"apiVersion\\\"\"; Kind string \"json:\\\"kind\\\"\" }"

// IsArgoCDApplication knows how to determine if a path in a filesystem is a valid path
func IsArgoCDApplication(fs *afero.Afero, targetPath string) (bool, error) {
	info, err := fs.Stat(targetPath)
	if err != nil {
		return false, fmt.Errorf("stating: %w", err)
	}

	if info.IsDir() {
		return false, nil
	}

	content, err := fs.ReadFile(targetPath)
	if err != nil {
		return false, fmt.Errorf("reading file: %w", err)
	}

	serializer := struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
	}{}

	err = yaml.Unmarshal(content, &serializer)
	if err != nil {
		if strings.Contains(err.Error(), argoCDUnmarshallingError) {
			return false, nil
		}

		return false, fmt.Errorf("unmarshalling: %w", err)
	}

	model := resources.GenerateDefaultArgoApp()

	if serializer.APIVersion != model.APIVersion {
		return false, nil
	}

	if serializer.Kind != model.Kind {
		return false, nil
	}

	return true, nil
}
