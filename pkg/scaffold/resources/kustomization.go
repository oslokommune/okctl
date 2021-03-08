package resources

// PatchTarget represents a target attribute in a PatchReference attribute
type PatchTarget struct {
	Kind string `json:"kind"`
}

// PatchReference represents an element in the patches attribute in a kustomization file
type PatchReference struct {
	Path   string      `json:"path"`
	Target PatchTarget `json:"target"`
}

// Kustomization represents a single kustomization.yaml file
type Kustomization struct {
	Resources []string         `json:"resources"`
	Patches   []PatchReference `json:"patches,omitempty"`
}

// AddPatch adds elements to the patches attribute of a kustomization file
func (k *Kustomization) AddPatch(ref PatchReference) {
	k.Patches = append(k.Patches, ref)
}

// AddResource adds elements to the resource attribute of a kustomization file
func (k *Kustomization) AddResource(resource string) {
	k.Resources = append(k.Resources, resource)
}

// NewKustomization initializes a Kustomization struct
func NewKustomization() *Kustomization {
	return &Kustomization{
		Resources: []string{},
		Patches:   []PatchReference{},
	}
}
