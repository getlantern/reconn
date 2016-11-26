package reconn

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"testing"
)

const (
	text = "hello world"
)

func TestReRead(t *testing.T) {
	doTest(t, 5, func(conn net.Conn) {
		rc := Wrap(conn, 5)
		b := make([]byte, len(text))

		// Read
		_, err := io.ReadFull(rc, b)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, text, string(b))

		b = make([]byte, 5)
		_, err = io.ReadFull(rc.Rereader(), b)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, text[:5], string(b))
	})
}

func doTest(t *testing.T, limit int, onConn func(net.Conn)) {
	l, err := net.Listen("tcp", ":0")
	if !assert.NoError(t, err) {
		return
	}
	defer l.Close()

	conn, err := net.Dial("tcp", l.Addr().String())
	if !assert.NoError(t, err) {
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(text))
	if !assert.NoError(t, err) {
		return
	}

	conn, err = l.Accept()
	if !assert.NoError(t, err) {
		return
	}
	defer conn.Close()

	onConn(conn)
}
