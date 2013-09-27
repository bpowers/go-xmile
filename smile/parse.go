// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smile

import (
	"fmt"
	"go/token"
	"strings"
	"unicode/utf8"
)

// Parse returns an abstract syntax tree corresponding to the given
// equation, or an error.
func Parse(name, eqn string) (Expr, error) {

	// it makes the lexer's code much cleaner to have a rune to
	// parse that marks the end of the equation
	if r, _ := utf8.DecodeLastRuneInString(eqn); r != ';' {
		eqn += ";"
	}

	fset := token.NewFileSet()
	f := fset.AddFile(name, fset.Base(), len(eqn))

	p := newParser(f, fset, newLexer(eqn, f))
	ast, ok := p.Parse()
	if p.errs.ErrorCount() != 0 {
		return nil, p.errs.GetErrorList(Sorted)
	} else if !ok {
		return nil, fmt.Errorf("p.Parse(): not ok %#v", p)
	}

	return ast, nil
}

type parser struct {
	tokf   *token.File
	fset   *token.FileSet
	lex    *lexer
	errs   ErrorVector
	levels []exprFn
}

func newParser(f *token.File, fs *token.FileSet, l *lexer) *parser {
	p := &parser{tokf: f, fset: fs, lex: l}
	p.levels = []exprFn{
		binaryLevelGen(0, p, "^"),
		binaryLevelGen(1, p, "*/"),
		binaryLevelGen(2, p, "+-"),
		p.factor,
	}
	return p
}

func (p *parser) expr() (Expr, bool) {
	return p.levels[0]()
}

func ident(tok Token) *Ident {
	return &Ident{tok.pos, tok.val}
}

func id(n string) *Ident {
	return &Ident{Name: n}
}

func (p *parser) Parse() (x Expr, ok bool) {
	if x, ok = p.expr(); !ok {
		return
	}
	la := p.lex.Peek()
	if la == nil {
		panic(fmt.Errorf("parser.Parse('%s'): missing semicolon", p.lex.s))
	}
	if la.kind != itemSemi {
		p.errorf(la, "expected end-of-equation, got %#v", la)
		return nil, false
	}
	p.lex.Token() // consume
	return
}

func (p *parser) errorf(tok *Token, f string, args ...interface{}) {
	var pos token.Position
	if tok != nil {
		pos = p.fset.Position(tok.pos)
	}
	p.errs.Error(pos, fmt.Sprintf(f, args...))
}

func floatLit(t *Token) *BasicLit {
	return &BasicLit{t.pos, token.FLOAT, t.val}
}

func opToken(t *Token) token.Token {
	switch t.val {
	case "^":
		return token.XOR // we interpret XOR as exponentiation
	case "+":
		return token.ADD
	case "-":
		return token.SUB
	case "*":
		return token.MUL
	case "/":
		return token.QUO
	}
	panic(fmt.Errorf("opToken(%#v): illegal token", t))
}

type exprFn func() (Expr, bool)

func binaryLevelGen(n int, p *parser, ops string) exprFn {
	return func() (lhs Expr, ok bool) {
		if p.lex.Peek() == nil {
			return nil, true
		}

		var next exprFn
		if n+1 >= len(p.levels) {
			panic(fmt.Errorf("binaryLevelGen(%d, '%s'): illegal level (max %d)",
				n, ops, len(p.levels)))
		}
		next = p.levels[n+1]

		if lhs, ok = next(); !ok {
			return
		}

		var op *Token
		for op, ok = p.consumeAnyOf(ops); ok; op, ok = p.consumeAnyOf(ops) {
			var rhs Expr
			if rhs, ok = next(); !ok {
				return
			}
			lhs = &BinaryExpr{
				X:     lhs,
				OpPos: op.pos,
				Op:    opToken(op),
				Y:     rhs,
			}
		}
		return lhs, true
	}
}

func (p *parser) factor() (x Expr, ok bool) {
	var lparen *Token
	if lparen, ok = p.consumeTok(itemLParen); ok {
		if x, ok = p.expr(); !ok {
			return
		}
		var rparen *Token
		if rparen, ok = p.consumeTok(itemRParen); !ok {
			p.errorf(p.lex.Peek(), "expected ')'")
			return nil, false
		}
		x = &ParenExpr{lparen.pos, x, rparen.pos}
		return
	}

	if x, ok = p.num(); ok {
		return
	} else if x, ok = p.ident(); ok {
		// CallExpr
		if tok, ok := p.consumeTok(itemLParen); ok {
			return p.call(x, tok)
		}
		return
	}

	p.errorf(p.lex.Peek(), "unexpected token")
	return nil, false
}

func (p *parser) call(fun Expr, lparen *Token) (x Expr, ok bool) {
	ce := &CallExpr{Fun: fun, Lparen: lparen.pos}
	x = ce

	var tok *Token
	if tok, ok = p.consumeTok(itemRParen); ok {
		ce.Rparen = tok.pos
		return
	}

	for {
		var arg Expr
		if arg, ok = p.expr(); !ok {
			p.errorf(p.lex.Peek(), "call: expected expr arg, not %#v", p.lex.Peek())
			return
		}
		ce.Args = append(ce.Args, arg)
		if _, ok = p.consumeAnyOf(","); ok {
			continue
		}
		if tok, ok = p.consumeTok(itemRParen); ok {
			ce.Rparen = tok.pos
			break
		}
		p.errorf(p.lex.Peek(), "call: expected ',' or ')', not %#v", p.lex.Peek())
		return nil, false
	}
	return
}

func (p *parser) ident() (Expr, bool) {
	if la := p.lex.Peek(); la != nil && la.kind == itemIdentifier {
		t := p.lex.Token()
		return &Ident{t.pos, t.val}, true
	}
	return nil, false
}

func (p *parser) num() (Expr, bool) {
	if la := p.lex.Peek(); la != nil && la.kind == itemNumber {
		t := p.lex.Token()
		return &BasicLit{t.pos, token.FLOAT, t.val}, true
	}
	return nil, false
}

func (p *parser) consumeAnyOf(ops string) (*Token, bool) {
	la := p.lex.Peek()
	if la == nil || la.kind != itemOperator {
		return nil, false
	}
	op, _ := utf8.DecodeRuneInString(la.val)
	if op != utf8.RuneError && strings.IndexRune(ops, op) > -1 {
		return p.lex.Token(), true
	}
	return nil, false
}

func (p *parser) consumeTok(ty itemType) (*Token, bool) {
	la := p.lex.Peek()
	if la == nil || la.kind != ty {
		return nil, false
	}
	return p.lex.Token(), true
}
