package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

var ErrEof = errors.New("EOF")

type TokenType int

const (
	TokenPlus   TokenType = iota // +
	TokenMinus                   // -
	TokenStar                    // *
	TokenDollar                  // $
	TokenNumber                  // number type. ex: double (1.23), integer (123)
	TokenString                  // string type.
)

var tokenTypeNames = []string{
	TokenPlus:   "PLUS",
	TokenMinus:  "MINUS",
	TokenStar:   "STAR",
	TokenDollar: "Dollar",
	TokenNumber: "Number",
	TokenString: "String",
}

func (t TokenType) String() string {
	return tokenTypeNames[t]
}

type Token struct {
	tokenType TokenType
	value     []byte
}

func (t Token) String() string {
	return fmt.Sprintf("Lexer{ tokenType:%d value:%+v}", t.tokenType, strconv.Quote(string(t.value)))
}

func newScanner(r io.Reader) *bufio.Scanner {
	var splitFn bufio.SplitFunc = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
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
	}

	scanner := bufio.NewScanner(r)
	scanner.Split(splitFn)

	return scanner
}

type Lexer struct {
	sc *bufio.Scanner
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		sc: newScanner(r),
	}
}

func (l *Lexer) Scan() bool {
	return l.sc.Scan()
}

func (l *Lexer) Tokens() []Token {
	bytes := l.sc.Bytes()
	tokens := []Token{}

	switch bytes[0] {
	case '+':
		tokens = append(tokens, Token{
			tokenType: TokenPlus,
			value:     bytes[0:1],
		})
	case '-':
		tokens = append(tokens, Token{
			tokenType: TokenMinus,
			value:     bytes[0:1],
		})
	case '*':
		tokens = append(tokens, Token{
			tokenType: TokenStar,
			value:     bytes[0:1],
		})
	case '$':
		tokens = append(tokens, Token{
			tokenType: TokenDollar,
			value:     bytes[0:1],
		})
	}

	if len(tokens) > 0 {
		prevToken := tokens[0]

		switch prevToken.tokenType {
		case TokenPlus:
		case TokenMinus:
			tokens = append(tokens, Token{
				tokenType: TokenString,
				value:     bytes[1:],
			})
		case TokenStar:
		case TokenDollar:
			tokens = append(tokens, Token{
				tokenType: TokenNumber,
				value:     bytes[1:],
			})
		default:
			tokens = append(tokens, Token{
				tokenType: TokenString,
				value:     bytes[1:],
			})
		}
	} else {
		tokens = append(tokens, Token{
			tokenType: TokenString,
			value:     bytes[:],
		})
	}

	return tokens
}

func (l *Lexer) Reset(r io.Reader) {
	if r == nil {
		l.sc = nil
	} else {
		l.sc = newScanner(r)
	}
}

func (l *Lexer) Err() error {
	return l.sc.Err()
}
