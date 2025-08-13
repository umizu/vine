package vine

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"syscall"
)

func parseRequest(fd int) *http.Request {
	var (
		req             http.Request
		read            int
		endOfHeadersIdx int
	)
	buf := make([]byte, 4*1024)

	for {
		n, err := syscall.Read(fd, buf[read:])
		if err != nil {
			slog.Error(err.Error(), "fd", fd)
		}
		read += n

		endOfHeadersIdx = bytes.Index(buf, []byte{'\r', '\n', '\r', '\n'}) // todo: not necessary to check entire buffer per read
		if endOfHeadersIdx != -1 {
			// todo: capture bytes from requestbody, if any
			break
		}
	}

	lines := strings.Split(string(buf[:endOfHeadersIdx]), "\r\n")
	fmt.Printf("start line: %s\n", lines[0]) // debugging

	return &req
}
