package kube

import "errors"

// ErrClusterNotFound represents the specified cluster is missing
var ErrClusterNotFound = errors.New("cluster not found")
