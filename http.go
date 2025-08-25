package vine

import "errors"

var ErrMissingHostHeader = errors.New("missing required Host header")

type Request struct {
	Headers map[string][]string
	Method  string
	Path    string

	Proto      string
	ProtoMajor int
	ProtoMinor int
}

func (r *Request) Validate() error {
	_, ok := r.Headers["Host"]
	if !ok {
		return ErrMissingHostHeader
	}
	return nil
}
