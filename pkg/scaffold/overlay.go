package scaffold

type applicationOverlay struct {
	Kustomization   []byte
	IngressPatch    []byte
	ServicePatch    []byte
	DeploymentPatch []byte
}

func newApplicationOverlay() applicationOverlay {
	return applicationOverlay{
		Kustomization:   []byte(""),
		IngressPatch:    []byte(""),
		ServicePatch:    []byte(""),
		DeploymentPatch: []byte(""),
	}
}
