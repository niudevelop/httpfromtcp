package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

// Parse consumes at most one header line (ending in CRLF) per call.
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlf := []byte("\r\n")
	i := bytes.Index(data, crlf)
	if i == -1 {
		return 0, false, nil // not enough data yet
	}

	// If CRLF is at the start, headers are done.
	if i == 0 {
		return 2, true, nil
	}

	line := string(data[:i]) // excludes CRLF

	colon := strings.IndexByte(line, ':')
	if colon <= 0 {
		return 0, false, fmt.Errorf("invalid header line")
	}

	rawKey := line[:colon]
	// Reject any whitespace around/in the key (catches "Host : ..." and " Host: ...")
	if strings.TrimSpace(rawKey) != rawKey || strings.IndexFunc(rawKey, func(r rune) bool {
		return r == ' ' || r == '\t'
	}) != -1 {
		return 0, false, fmt.Errorf("invalid header key spacing")
	}

	rawVal := line[colon+1:]
	key := rawKey
	val := strings.TrimSpace(rawVal)

	h[key] = val

	// consumed: header line + CRLF
	return i + 2, false, nil
}
