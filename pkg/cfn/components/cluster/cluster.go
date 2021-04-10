// Package cluster knows how to return a consistent cluster name
package cluster

// Cluster provides a type for
// creating a consistent name
type Cluster struct {
	StoredName string
}

// Name returns the name of the cluster
func (c *Cluster) Name() string {
	return c.StoredName
}

// New returns a structure for printing
// a consistent Cluster name for use in a
// cloud formation template
func New(clusterName string) *Cluster {
	return &Cluster{
		StoredName: clusterName,
	}
}
