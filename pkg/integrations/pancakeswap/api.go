package pancakeswap

import "errors"

type v2ApiResponse[Data any] struct {
	UpdatedAt int64 `json:"updated_at"`
	Data      Data  `json:"data"`
}

var ErrRequestInvalid = errors.New("request invalid")
var ErrRequestFailed = errors.New("request failed")
var ErrResponseInvalid = errors.New("response invalid")
