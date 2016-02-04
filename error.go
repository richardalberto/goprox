package goprox

import "encoding/json"

type Error struct {
	Message error `json:"message"`
}

func NewError(msg error) *Error {
	return &Error{Message: msg}
}

func (e *Error) JSON() ([]byte, error) {
	return json.Marshal(e)
}
