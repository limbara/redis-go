package resp

import (
	"io"

	"github.com/codecrafters-io/redis-starter-go/app/lexer"
)

type CommandParser struct {
	lexer   *lexer.Lexer
	tokens  []lexer.Token
	isLexed bool
}

func NewCommandParser(r io.Reader) *CommandParser {
	newLexer := lexer.NewLexer(r)

	return &CommandParser{
		lexer:   newLexer\,
		tokens:  make([]lexer.Token, 0),
		isLexed: false,
	}
}
