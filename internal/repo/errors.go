package repo

import "errors"

var ErrNotFound = errors.New("item not found")
var ErrConflict = errors.New("id already exists")
var ErrInvalidEmail = errors.New("invalid email")
