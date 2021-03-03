package scaffold

type applicationOverlay struct {
	IngressPatch    []byte
	ServicePatch    []byte
	DeploymentPatch []byte
}

func newApplicationOverlay() applicationOverlay {
	return applicationOverlay{
		IngressPatch:    []byte(""),
		ServicePatch:    []byte(""),
		DeploymentPatch: []byte(""),
	}
}
