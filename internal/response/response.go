package response

import (
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"strconv"

	"github.com/niudevelop/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCode200 StatusCode = 200
	StatusCode400 StatusCode = 400
	StatusCode500 StatusCode = 500
)
const crlf = "\r\n"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusCode200:

		b := []byte(fmt.Sprintf("HTTP/1.1 %d %s%s", statusCode, "OK", crlf))
		if _, err := w.Write(b); err != nil {
			return err
		}
	case StatusCode400:
		b := []byte(fmt.Sprintf("HTTP/1.1 %d %s%s", statusCode, "Bad Request", crlf))
		if _, err := w.Write(b); err != nil {
			return err
		}
	case StatusCode500:
		b := []byte(fmt.Sprintf("HTTP/1.1 %d %s%s", statusCode, "Internal Server Error", crlf))
		if _, err := w.Write(b); err != nil {
			return err
		}
	default:
		return errors.New("unknown status code")
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.NewHeaders()
	header.Set("Content-Length", strconv.Itoa(contentLen))
	header.Set("Connection", "close")
	header.Set("Content-Type", "text/plain")
	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for k, v := range headers {
		if _, err := w.Write([]byte(fmt.Sprintf("%s: %s%s", canonicalHeaderKey(k), v, crlf))); err != nil {
			return err
		}
	}
	return nil
}
func canonicalHeaderKey(s string) string {
	return textproto.CanonicalMIMEHeaderKey(s)
}
