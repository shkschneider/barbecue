package core

import (
	"errors"
)

var (
	ErrNotFound = errors.New("NotFound")
	ErrNothing = errors.New("Nothing")
	ErrAmbiguous = errors.New("TooAmbiguous")
	ErrApi = errors.New("Api")
	ErrDatabase = errors.New("Database")
	ErrUnknown = errors.New("Unknown")
)
