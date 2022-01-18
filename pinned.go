// This provides us with a way of pinning the k8s cli's to the
// same version without breaking go fmt et al
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
//go:build pinned
// +build pinned

package okctl

import (
	_ "github.com/containerd/containerd"
	_ "github.com/docker/distribution"
)
