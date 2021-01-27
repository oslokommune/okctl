package okctl

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/state"
)

// ErrorEnvironmentNotFound indicates that the requested environment isn't available
type ErrorEnvironmentNotFound struct {
	TargetEnvironment     string
	AvailableEnvironments []string
}

func (e ErrorEnvironmentNotFound) Error() string {
	return fmt.Sprintf("\"%s\" is not in available environments %v", e.TargetEnvironment, e.AvailableEnvironments)
}

func getEnvironments(clusters map[string]state.Cluster) []string {
	environments := make([]string, len(clusters))

	for env := range clusters {
		environments = append(environments, env)
	}

	return environments
}
