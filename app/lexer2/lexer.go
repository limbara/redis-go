package lexer2

import (
	"bufio"
	"fmt"
	"strconv"
)

type TokenType int

const (
	TokenError  TokenType = iota
	TokenPlus             // +
	TokenMinus            // -
	TokenStar             // *
	TokenDollar           // $
	TokenNumber           // number type. ex: double (1.23), integer (123)
	TokenString           // string type.
)

var tokenTypeNames = []string{
	TokenError:  "Error",
	TokenPlus:   "Plus",
	TokenMinus:  "Minus",
	TokenStar:   "Star",
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
	return fmt.Sprintf("Token{ tokenType:%s value:%+v}", t.tokenType, strconv.Quote(string(t.value)))
}

type Lexer struct {
	br *bufio.Reader
}
