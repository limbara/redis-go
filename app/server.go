package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
)

var ErrEof = errors.New("EOF")

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err, ok := recover().(error); err != nil && ok {
			fmt.Println("recover error handle connection: ", err)

			switch {
			// close connection if broken pipe or [ErrEof]
			case errors.Is(err, ErrEof), errors.Is(err, syscall.EPIPE):
				conn.Close()
			}
		}
	}()

	for {
		scanner := bufio.NewScanner(conn)
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			// conn read hit ioEOF if connection closed by client
			if atEOF {
				return 0, nil, ErrEof
			}

			advance, token, err = bufio.ScanLines(data, atEOF)
			isEOF := len(data)-advance == 0

			if isEOF && len(token) != 0 {
				return advance, token, bufio.ErrFinalToken
			} else {
				return advance, token, err
			}
		})

		for scanner.Scan() {
			fmt.Println("Scanned", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			if opError, ok := err.(*net.OpError); ok {
				fmt.Printf("It's op error: %+v\n", opError)
				panic(opError.Err)
			}

			panic(err)
		}

		fmt.Println("Writing")
		_, err := io.WriteString(conn, "+PONG\r\n")
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
