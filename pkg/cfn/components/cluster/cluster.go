package cluster

import "fmt"

type cluster struct {
	name string
}

func (c *cluster) Name() string {
	return c.name
}

func New(name, env string) *cluster {
	return &cluster{
		name: fmt.Sprintf("%s-%s", name, env),
	}
}
