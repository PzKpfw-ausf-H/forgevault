package minio

import "errors"

var ErrMissingEndpoint = errors.New("missing endpoint")
var ErrMissingAccess = errors.New("missing access key")
var ErrMissingSecret = errors.New("missing secret key")
var ErrMissingSSL = errors.New("missing / wrong SSL")
