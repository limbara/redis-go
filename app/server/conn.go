package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/app/lexer"
)

var lexerReaderPool sync.Pool

func newLexer(r io.Reader) *lexer.Lexer {
	if v := lexerReaderPool.Get(); v != nil {
		lexer := v.(*lexer.Lexer)
		lexer.Reset(r)

		return lexer
	}

	return lexer.NewLexer(r)
}

func putLexer(l *lexer.Lexer) {
	l.Reset(nil)
	lexerReaderPool.Put(l)
}

type conn struct {
	// rwc is the underlying network connection.
	// This is never wrapped by other types and is the value given out
	// to CloseNotifier callers. It is usually of type *net.TCPConn orlim
	// *tls.Conn.
	rwc net.Conn

	lexer *lexer.Lexer
}

func NewConn(rwc net.Conn) *conn {
	return &conn{
		rwc: rwc,
	}
}

func (c *conn) Serve() {
	defer func() {
		if err, ok := recover().(error); err != nil && ok {
			fmt.Println("recover error handle connection: ", err)

			switch {
			// close connection if broken pipe or [ErrEof]
			case errors.Is(err, lexer.ErrEof), errors.Is(err, syscall.EPIPE):
				c.close()
			}
		}
	}()

	c.lexer = newLexer(c.rwc)
	for {
		c.lexer.Reset(c.rwc)

		for c.lexer.Scan() {
			fmt.Printf("Scanned: %+v\n", c.lexer.Tokens())
		}

		if err := c.lexer.Err(); err != nil {
			if opError, ok := err.(*net.OpError); ok {
				fmt.Printf("It's op error: %+v\n", opError)
				panic(opError.Err)
			}

			panic(err)
		}

		fmt.Println("Writing")
		_, err := io.WriteString(c.rwc, "+PONG\r\n")
		if err != nil {
			if opError, ok := err.(*net.OpError); ok {
				fmt.Printf("It's op error: %+v\n", opError)
				panic(opError.Err)
			}
			panic(err)
		}
		fmt.Println()
	}
}

func (c *conn) flush() {
	if c.lexer != nil {
		putLexer(c.lexer)
		c.lexer = nil
	}
}

func (c *conn) close() {
	c.flush()
	c.rwc.Close()
}
