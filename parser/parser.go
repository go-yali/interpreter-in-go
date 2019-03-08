package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l *lexer.Lexer // pointer to an instance of the lexer

	// similar to 'pointers' in our lexer (position and readPosition)
	// But instead of pointing to a charcter of the input, they point to the current and next token
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokents, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

// To get the next tokens
// Advances both curToken and peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
