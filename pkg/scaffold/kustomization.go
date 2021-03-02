package scaffold

type PatchTarget struct {
	Kind string `json:"kind"`
}

type PatchReference struct {
	Path   string      `json:"path"`
	Target PatchTarget `json:"target"`
}

type Kustomization struct {
	Resources []string         `json:"resources"`
	Patches   []PatchReference `json:"patches"`
}

func (k *Kustomization) AddPatch(ref PatchReference) {
	k.Patches = append(k.Patches, ref)
}

func (k *Kustomization) AddResource(resource string) {
	k.Resources = append(k.Resources, resource)
}

func NewKustomization() *Kustomization {
	return &Kustomization{
		Resources: []string{},
		Patches:   []PatchReference{},
	}
}
