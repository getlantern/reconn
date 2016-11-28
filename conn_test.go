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

func TestReReadOkay(t *testing.T) {
	doTest(t, 5, func(conn net.Conn) {
		rc := Wrap(conn, 5)
		b := make([]byte, 5)

		// Read
		_, err := io.ReadFull(rc, b)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, text[:5], string(b))

		rr, err := rc.Rereader()
		if !assert.NoError(t, err) {
			return
		}
		b = make([]byte, 5)
		_, err = io.ReadFull(rr, b)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, text[:5], string(b))
	})
}

func TestReReadOverflow(t *testing.T) {
	doTest(t, 5, func(conn net.Conn) {
		rc := Wrap(conn, 5)
		b := make([]byte, len(text))

		// Read
		_, err := io.ReadFull(rc, b)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, text, string(b))

		_, err = rc.Rereader()
		assert.Equal(t, ErrOverflowed, err)
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
