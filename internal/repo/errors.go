package repo

import "errors"

var (
	ErrNotFound      = errors.New("repository: not found")
	ErrAlreadyExists = errors.New("repository: already exists")
)
