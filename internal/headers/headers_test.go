package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	t.Run("Valid single header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "localhost:42069", headers["Host"])
		assert.Equal(t, 23, n)
		assert.False(t, done)
	})

	t.Run("Valid single header with extra whitespace", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host:    localhost:42069       \r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		assert.Equal(t, "localhost:42069", headers["Host"])
		assert.Equal(t, len("Host:    localhost:42069       \r\n"), n)
		assert.False(t, done)
	})

	t.Run("Valid 2 headers with existing headers", func(t *testing.T) {
		headers := NewHeaders()
		headers["Host"] = "localhost:42069"

		data1 := []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")

		n, done, err := headers.Parse(data1)
		require.NoError(t, err)
		assert.False(t, done)
		assert.Equal(t, "curl/7.81.0", headers["User-Agent"])
		assert.Equal(t, len("User-Agent: curl/7.81.0\r\n"), n)

		data2 := data1[n:]
		n2, done2, err2 := headers.Parse(data2)
		require.NoError(t, err2)
		assert.False(t, done2)
		assert.Equal(t, "*/*", headers["Accept"])
		assert.Equal(t, len("Accept: */*\r\n"), n2)
	})

	t.Run("Valid done", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		assert.Equal(t, 2, n)
		assert.True(t, done)
	})

	t.Run("Invalid spacing header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("       Host : localhost:42069       \r\n\r\n")

		n, done, err := headers.Parse(data)

		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})
}
