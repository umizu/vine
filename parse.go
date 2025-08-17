package vine

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	ErrMalformedReq = errors.New("invalid request")
)

func parseRequest(req *Request, conn net.Conn) error {
	var startLineParsed bool
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if len(line) > 1 && line[len(line)-2] == '\r' { // remove CRLF / LF
			line = line[:len(line)-2]
		} else {
			line = line[:len(line)-1]
		}

		// start-line
		if !startLineParsed {
			if err := parseRequestLine(req, line); err != nil {
				return err
			}
			startLineParsed = true
			continue
		}

		// headers
		if line == "" { // request end
			return nil
		}

		if line[0] == ' ' || line[0] == '\t' {
			continue
		}

		colonIdx := strings.IndexByte(line, ':')
		if colonIdx == -1 {
			return fmt.Errorf("malformed field-line: %q", line)
		}
		hKey := toPascalCaseHeader(line[:colonIdx])
		hVal := strings.TrimSpace(line[colonIdx+1:])

		if req.Headers == nil {
			req.Headers = make(map[string][]string)
		}

		req.Headers[hKey] = strings.Split(hVal, ",")
	}

}

func parseRequestLine(r *Request, line string) error {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return ErrMalformedReq
	}

	r.Method = parts[0] // todo: validate?
	r.Path = parts[1]   // todo: validate?
	r.Proto = parts[2]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 || versionParts[0] != "HTTP" {
		return ErrMalformedReq
	}

	protoParts := strings.Split(versionParts[1], ".")
	if len(protoParts) != 2 {
		return ErrMalformedReq
	}

	protoMajor, err := strconv.Atoi(protoParts[0])
	if err != nil {
		return ErrMalformedReq
	}
	r.ProtoMajor = protoMajor

	protoMinor, err := strconv.Atoi(protoParts[1])
	if err != nil {
		fmt.Print("hit3")
		return ErrMalformedReq
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
