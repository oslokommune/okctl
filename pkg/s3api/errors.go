package s3api

import "errors"

// ErrBucketDoesNotExist represents that the specified bucket does not exist
var ErrBucketDoesNotExist = errors.New("bucket does not exist")
