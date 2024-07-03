package lexer2

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

const minBufferSize = 16
const readBufferSize = 64
const maxConsecutiveEmptyReads = 100

var errNegativeRead = errors.New("LexReader: reader returned negative count from Read")

type LexReader struct {
	rd     io.Reader
	err    error
	buffer []byte
	r, w   int // buf read and write positions
	p      int // the current position of out scan
}

func NewReader(r io.Reader) *LexReader {
	lr := new(LexReader)
	lr.reset(make([]byte, minBufferSize), r)

	return lr
}

func (lr *LexReader) fill() {
	// Read new data: try a limited number of times.
	for i := maxConsecutiveEmptyReads; i > 0; i-- {
		// p := make([]byte, readBufferSize)
		n, err := lr.rd.Read(lr.buffer[lr.w:])
		fmt.Println("Readed\t", strconv.Quote(string(lr.buffer)), "\nLength:", n)
		if n < 0 {
			panic(errNegativeRead)
		}

		if err != nil {
			lr.err = err
			lr.w += n
			fmt.Println("err != nil\t", lr.buffer)
			return
		}

		if n > 0 {
			lr.w += n
			fmt.Println("err == nil\t", lr.buffer)
			return
		}
	}
	lr.err = io.ErrNoProgress
}

func (lr *LexReader) Next() (rune, error) {
	for lr.r+utf8.UTFMax > lr.w && !utf8.FullRune(lr.buffer[lr.r:lr.w]) && lr.err == nil {
		lr.fill()
	}

	if lr.r == lr.w {
		lr.err = io.EOF

		return utf8.RuneError, lr.err
	}

	r, size := rune(lr.buffer[lr.p]), 1
	if r >= utf8.RuneSelf {
		r, size = utf8.DecodeRune(lr.buffer[lr.p:lr.w])
	}
	lr.p += size

	return r, nil
}

func (lr *LexReader) Backup() {
	if lr.err != io.EOF && lr.p > 0 {
		_, w := utf8.DecodeLastRune(lr.buffer[:lr.p])
		lr.p -= w

		if lr.p < lr.r { // if
			lr.r = lr.p
		}
	}
}

func (lr *LexReader) Peek() (rune, error) {
	r, e := lr.Next()
	lr.Backup()

	return r, e
}

func (lr *LexReader) Ignore() {
	lr.r = lr.p
}

func (lr *LexReader) Bytes() []byte {
	bytes := lr.buffer[lr.r:lr.p]

	lr.r = lr.p

	return bytes
}

func (lr *LexReader) Err() error {
	return lr.err
}

func (lr *LexReader) Reset(r io.Reader) {
	lr.reset(make([]byte, minBufferSize), r)
}

func (b *LexReader) reset(buf []byte, r io.Reader) {
	*b = LexReader{
		buffer: buf,
		rd:     r,
	}
}
