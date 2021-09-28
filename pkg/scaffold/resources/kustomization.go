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
	k.RemovePatch(ref.Path)

	k.Patches = append(k.Patches, ref)
}

// AddResource adds elements to the resource attribute of a kustomization file
func (k *Kustomization) AddResource(resource string) {
	k.RemoveResource(resource)

	k.Resources = append(k.Resources, resource)
}

// RemoveResource removes a resource
func (k *Kustomization) RemoveResource(path string) {
	index := k.resourceIndex(path)

	if index == -1 {
		return
	}

	k.Resources[index] = k.Resources[len(k.Resources)-1]
	k.Resources = k.Resources[:len(k.Resources)-1]
}

// RemovePatch removes a PatchReference
func (k *Kustomization) RemovePatch(path string) {
	index := k.patchIndex(PatchReference{Path: path})

	if index == -1 {
		return
	}

	k.Patches[index] = k.Patches[len(k.Patches)-1]
	k.Patches = k.Patches[:len(k.Patches)-1]
}

func (k *Kustomization) resourceIndex(resource string) int {
	for index, item := range k.Resources {
		if item == resource {
			return index
		}
	}

	return -1
}

func (k *Kustomization) patchIndex(ref PatchReference) int {
	for index, item := range k.Patches {
		if item.Path == ref.Path {
			return index
		}
	}

	return -1
}

// NewKustomization initializes a Kustomization struct
func NewKustomization() *Kustomization {
	return &Kustomization{
		Resources: []string{},
		Patches:   []PatchReference{},
	}
}
