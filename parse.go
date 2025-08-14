package vine

import (
	"bytes"
	"fmt"
	"log/slog"
	"net"
	"strings"
)

func parseRequest(conn net.Conn) *Request {
	var (
		req             Request
		read            int
		endOfHeadersIdx int
	)
	buf := make([]byte, 4*1024)

	for {
		n, err := conn.Read(buf[read:])
		if err != nil {
			slog.Error("conn read", "err", err)
		}
		read += n

		endOfHeadersIdx = bytes.Index(buf, []byte{'\r', '\n', '\r', '\n'}) // todo: not necessary to check entire buffer per read
		if endOfHeadersIdx != -1 {
			// todo: capture bytes from request body, if any
			break
		}
	}

	lines := strings.Split(string(buf[:endOfHeadersIdx]), "\r\n")
	fmt.Printf("start line: %s\n", lines[0]) // debugging

	return &req
}
