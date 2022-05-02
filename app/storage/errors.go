package storage

import "errors"

var ErrWriting = errors.New("error writing")
var ErrReading = errors.New("error reading")
var ErrNotFound = errors.New("not found")
