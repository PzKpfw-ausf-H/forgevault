package minio

import "errors"

var ErrMissingEndpoint = errors.New("mising endpoint")
var ErrMissingAccess = errors.New("missing access key")
var ErrMissingSecret = errors.New("mising secret key")
