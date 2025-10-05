package inmem

import "errors"

var ErrNotFound = errors.New("entity not found")
var ErrAlreadyExists = errors.New("entity already exists")
