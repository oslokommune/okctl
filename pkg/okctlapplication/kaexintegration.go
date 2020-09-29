package okctlapplication

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/oslokommune/kaex/pkg/api"
)

/*
AcquireApplication returns an okctl Application based on stdin or a file
*/
func AcquireApplication(path string) (api.Application, error) {
	var (
		rawApplication []byte
		err            error
	)

	if path == "-" {
		rawApplication, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return api.Application{}, fmt.Errorf("failed to read stdin: %w", err)
		}
	} else {
		rawApplication, err = ioutil.ReadFile(filepath.Clean(path))
		if err != nil {
			return api.Application{}, fmt.Errorf("failed to read file: %w", err)
		}
	}

	app, err := api.ParseApplication(string(rawApplication))
	if err != nil {
		return api.Application{}, fmt.Errorf("unable to parse application: %w", err)
	}

	if app.Ingress.Annotations == nil {
		app.Ingress.Annotations = map[string]string{}
	}

	app.Ingress.Annotations["kubernetes.io/ingress.class"] = "alb"
	app.Ingress.Annotations["alb.ingress.kubernetes.io/scheme"] = "internet-facing"

	return app, err
}

/*
ArgoCDDeploymentResources contains the necessary resources for a ArgoCD deployment
*/
type ArgoCDDeploymentResources struct {
	KubernetesResourcesBuffer bytes.Buffer
	ArgoAppBuffer             bytes.Buffer
}

/*
ConvertApplicationToResources turns an api.Application into a ArgoCDDeploymentResources containing buffers with relevant
*/
func ConvertApplicationToResources(app api.Application, iacRepoURL string) (*ArgoCDDeploymentResources, error) {
	if iacRepoURL == "" {
		iacRepoURL = "git@github.com:<organization>/<repository>"
	}

	expandedApp := ArgoCDDeploymentResources{}

	err := api.Expand(&expandedApp.KubernetesResourcesBuffer, app, false)
	if err != nil {
		return &expandedApp, fmt.Errorf("error expanding application %w", err)
	}

	argoApp, err := CreateArgoApp(app, iacRepoURL)
	if err != nil {
		return &expandedApp, fmt.Errorf("error creating ArgoApp from application.yaml: %w", err)
	}

	err = api.WriteResource(&expandedApp.ArgoAppBuffer, argoApp)
	if err != nil {
		return &expandedApp, fmt.Errorf("error writing ArgoApp to buffer: %w", err)
	}

	return &expandedApp, nil
}
