package vine

type Request struct {
	Headers map[string][]string
	Method  string
	Path    string

	Proto      string
	ProtoMajor int
	ProtoMinor int
}
