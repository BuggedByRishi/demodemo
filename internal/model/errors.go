package model

import "errors"

// ErrNotFound is returned by repositories when a record doesn't exist.
var ErrNotFound = errors.New("record not found")

// ErrConflict is returned when a unique constraint would be violated.
var ErrConflict = errors.New("record already exists")

