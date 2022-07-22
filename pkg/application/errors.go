package application

import "errors"

var ErrReading = errors.New("error reading")
var ErrWriting = errors.New("error writing")

var ErrNotFound = errors.New("not found")
