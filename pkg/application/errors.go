package application

import "errors"

var ErrReading = errors.New("error reading")
var ErrWriting = errors.New("error writing")

var ErrInvalid = errors.New("error invalid")
var ErrNotFound = errors.New("error not found")
