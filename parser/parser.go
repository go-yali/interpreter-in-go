package parser

import (
	"fmt"
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

	errors []string

	// allows us to check if the appropriate map has a parsing function associated with curToken.Type
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	// prefixParseFn gets called when we encounter the associated token type in prefix position
	prefixParseFn func() ast.Expression

	// infixParseFn gets called when we encounter the token type in infix position
	infixParseFn func(ast.Expression) ast.Expression // takes the 'left side' of the infix operator
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

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

	//construct the root node of the AST
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// iterates (by repeatedly calling nextToken) over every token in the input until it encounters an EOF
	for !p.curTokenIs(token.EOF) {

		// during each iteration it parses a statement
		stmt := p.parseStatement()

		// unless the statement is nil, it adds the statement to the program's list of statements
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	// Finally, the root node is returned
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// parseLetStatement constructs an *ast.LetStatement node with the token its currently sitting on (a LET token), then advances the tokens while making assertions about the next token with calls to expectPeek
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// First, an Identifier is expected
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Then, an equal sign is expected
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we encounter a semicolon

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: We're skipping the expressions until we encounter a SEMICOLON

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek gets called, for instance during a Let Statement parsing, when we expect the next (peek) token to be something specific : Let statements have an name or identifier, an assign token, and then a value/expression
// Only if it is correct does expectPeek advance the tokens by calling nextToken
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// helper methods that add entries to the prefixParseFns & infixParseFns maps
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
