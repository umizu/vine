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
		read            int
		startlineEndIdx int
		startLineParsed bool
	)

	buf := make([]byte, 4*1024)
	for {
		n, err := conn.Read(buf[read:])
		if err != nil {
			slog.Error("conn read", "err", err)
		}
		read += n

		if !startLineParsed {
			startlineEndIdx = bytes.Index(buf[:read], []byte{'\r', '\n'})
			if startlineEndIdx == -1 {
				continue
			}
			if err := parseStartLine(&req, string(buf[:startlineEndIdx])); err != nil {
				return nil, err
			}
			startLineParsed = true
			fmt.Print("startLine parsed successfully")
		}

		break
	}

	return &req, nil
}

func parseStartLine(r *Request, line string) error {
	fmt.Println("parsing start line:", line)

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
