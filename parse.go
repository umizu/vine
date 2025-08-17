package vine

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
)

var (
	ErrInvalidReqFormat = errors.New("invalid request")
)

func parseRequest(conn net.Conn) (*Request, error) {
	var (
		req             Request
		totalRead       int
		lastLineEndIdx  int
		startLineParsed bool
	)
	buf := make([]byte, 4*1024)

	for {
		n, err := conn.Read(buf[totalRead:])
		if err != nil {
			slog.Error("conn read", "err", err)
		}
		totalRead += n

		// start line
		if !startLineParsed {
			lineEnd := bytes.Index(buf[:totalRead], []byte{'\r', '\n'})
			if lineEnd == -1 {
				continue
			}
			if err := parseStartLine(&req, string(buf[:lineEnd])); err != nil {
				return nil, err
			}
			startLineParsed = true
			lastLineEndIdx = lineEnd + 2
		}

		// headers
		for {
			lineEnd := bytes.Index(buf[lastLineEndIdx:totalRead], []byte("\r\n"))
			if lineEnd == -1 {
				break
			}

			lineEnd += lastLineEndIdx
			line := string(buf[lastLineEndIdx:lineEnd])
			lastLineEndIdx = lineEnd + 2

			if line == "" { // request end
				return &req, nil
			}

			colonIdx := strings.IndexByte(line, ':')
			if colonIdx == -1 {
				return nil, fmt.Errorf("malformed header: %q", line)
			}
			hKey := toPascalCaseHeader(line[:colonIdx])
			hVal := strings.TrimSpace(line[colonIdx+1:])

			if req.Headers == nil {
				req.Headers = make(map[string][]string)
			}

			req.Headers[hKey] = strings.Split(hVal, ",")
		}
	}
}

func parseStartLine(r *Request, line string) error {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return ErrInvalidReqFormat
	}

	r.Method = parts[0] // todo: validate?
	r.Path = parts[1]   // todo: validate?
	r.Proto = parts[2]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 || versionParts[0] != "HTTP" {
		return ErrInvalidReqFormat
	}

	protoParts := strings.Split(versionParts[1], ".")
	if len(protoParts) != 2 {
		return ErrInvalidReqFormat
	}

	protoMajor, err := strconv.Atoi(protoParts[0])
	if err != nil {
		return ErrInvalidReqFormat
	}
	r.ProtoMajor = protoMajor

	protoMinor, err := strconv.Atoi(protoParts[1])
	if err != nil {
		return ErrInvalidReqFormat
	}
	r.ProtoMinor = protoMinor

	return nil
}

func toPascalCaseHeader(header string) string {
	parts := strings.Split(header, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
		}
	}
	return strings.Join(parts, "-")
}
