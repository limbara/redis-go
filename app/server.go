package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

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
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error accepting connection: ", err)
			conn.Close()
		}
	}()

	for {
		scanner := bufio.NewScanner(conn)
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			advance, token, err = bufio.ScanLines(data, atEOF)
			isEOF := len(data)-advance == 0
			if isEOF {
				return advance, token, bufio.ErrFinalToken
			} else {
				return advance, token, err
			}
		})

		for scanner.Err() == nil && scanner.Scan() {
			fmt.Println("Scanned", scanner.Text())
		}
		if scanner.Err() != nil {
			fmt.Println("Error Scanner", scanner.Err())
		}

		fmt.Println("Writing")
		_, err := io.WriteString(conn, "+PONG\r\n")
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
