package cluster

import "fmt"

type Cluster struct {
	StoredName string
}

func (c *Cluster) Name() string {
	return c.StoredName
}

// New returns a structure for printing
// a consistent Cluster name for use in a
// cloud formation template
func New(name, env string) *Cluster {
	return &Cluster{
		StoredName: fmt.Sprintf("%s-%s", name, env),
	}
}
