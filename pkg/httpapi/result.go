package httpapi

type Response struct {
	Status  string      `json:"status"`
	Payload interface{} `json:"payload,omitempty"`
}

func NewResponse(status string) *Response {
	return &Response{Status: status}
}
func NewOkResponse() *Response {
	return NewResponse("ok")
}
func NewErrorResponse() *Response {
	return NewResponse("error")
}
func (r *Response) WithPayload(payload interface{}) *Response {
	r.Payload = payload
	return r
}
