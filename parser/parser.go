package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type (
	// prefixParseFn gets called when we encounter the associated token type in prefix position
	prefixParseFn func() ast.Expression

	// infixParseFn gets called when we encounter the token type in infix position
	infixParseFn func(ast.Expression) ast.Expression // takes the 'left side' of the infix operator
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

// User iota to increment these constants starting at 1 for LOWEST and 7 for CALL
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Initialize the prefixParseFns map on Parser and register a parsing function. Do the same for infixParseFns
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	// If, for eg, we encounter a token in a prefix expression
	// of type: token.IDENT, the parsing function
	// to call is parseIdentifier. Same for infix expressions
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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

// Main idea of Pratt parser: association of parsing functions with token types. EG: When I encounter LET token type, appropriate parseLetStatement() function is called
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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

// parseExpressionStatement constructs an AST node, and only advance curToken if the next token is a semicolon
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// All parsing functions, this one, prefixParseFun, and infixParseFn - don't advance tokens.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {

	// Check: Do we have a parsing function associated with p.curToken.Type in the prefix position?
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// parses the literal "5" from input into the numeric expression
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// Builds an AST node, like usual
// BUT: It advances our tokens by calling p.nextToken()!
// (Because We're working with prefix and expression)
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()

	// Now, when parseExpression is called, tokens have been advanced
	// in the case of `-5`, below: p.curToken.Type is token.INT.
	// parseExpression then checks the registered prefix parsing functions,
	// finds parseIntegerLiteral,
	// which builds an *ast.IntegerLiteral node and returns it.
	// parseExpression returns this newly constructed node and parsePrefixExpression uses it to fill the Right field of *ast.PrefixExpression.

	expression.Right = p.parseExpression(PREFIX)
	return expression
}

// Every time we make a new parser, we add an Expression for each operator token
// parseInfixExpression:
// 1. Takes argument left expression
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {

	// 2. constructs an InfixExpression node
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	// 3. assigns the precedence of the current token (which is the infix operator) to local var precedence
	precedence := p.curPrecedence()
	// 4. advances tokens
	p.nextToken()
	// 5. fills in expression.Right with another call to parseExpression
	expression.Right = p.parseExpression(precedence)
	return expression
}

//// HELPER METHODS ////

// helper methods that add entries to the prefixParseFns & infixParseFns maps
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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

// peekPrecedence returns the precedence associate with the token type of p.peekToken, defaulting to the lowest
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curPrecedence returns the precedence associate with the token type of p.peekToken, defaulting to the lowest
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse functions for %s found", t)
	p.errors = append(p.errors, msg)
}
