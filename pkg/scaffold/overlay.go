package scaffold

// ApplicationOverlay contains the contents of the kustomize patches
type ApplicationOverlay struct {
	Kustomization   []byte
	IngressPatch    []byte
	ServicePatch    []byte
	DeploymentPatch []byte
}

func newApplicationOverlay() ApplicationOverlay {
	return ApplicationOverlay{
		Kustomization:   []byte(""),
		IngressPatch:    []byte(""),
		ServicePatch:    []byte(""),
		DeploymentPatch: []byte(""),
	}
}
