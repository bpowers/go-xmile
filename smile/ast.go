// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smile

import (
	"go/token"
)

type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

type Expr interface {
	Node
	exprNode()
}

// An expression is represented by a tree consisting of one
// or more of the following concrete expression nodes.
//
type (
	// A BadExpr node is a placeholder for expressions containing
	// syntax errors for which no correct expression nodes can be
	// created.
	//
	BadExpr struct {
		From, To token.Pos // position range of bad expression
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
	}
	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		ValuePos token.Pos   // literal position
		Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
		Value    string      // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
	}

	// A ParenExpr node represents a parenthesized expression.
	ParenExpr struct {
		Lparen token.Pos // position of "("
		X      Expr      // parenthesized expression
		Rparen token.Pos // position of ")"
	}

	// An IndexExpr node represents an expression followed by an index.
	IndexExpr struct {
		X      Expr      // expression
		Lbrack token.Pos // position of "["
		Index  Expr      // index expression
		Rbrack token.Pos // position of "]"
	}

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr struct {
		Fun    Expr      // function expression
		Lparen token.Pos // position of "("
		Args   []Expr    // function arguments; or nil
		Rparen token.Pos // position of ")"
	}

	// A UnaryExpr node represents a unary expression.
	// Unary "*" expressions are represented via StarExpr nodes.
	//
	UnaryExpr struct {
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		X     Expr        // operand
	}

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X     Expr        // left operand
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		Y     Expr        // right operand
	}
)

func (x *BadExpr) Pos() token.Pos    { return x.From }
func (x *Ident) Pos() token.Pos      { return x.NamePos }
func (x *BasicLit) Pos() token.Pos   { return x.ValuePos }
func (x *ParenExpr) Pos() token.Pos  { return x.Lparen }
func (x *IndexExpr) Pos() token.Pos  { return x.X.Pos() }
func (x *CallExpr) Pos() token.Pos   { return x.Fun.Pos() }
func (x *UnaryExpr) Pos() token.Pos  { return x.OpPos }
func (x *BinaryExpr) Pos() token.Pos { return x.X.Pos() }

func (x *BadExpr) End() token.Pos    { return x.To }
func (x *Ident) End() token.Pos      { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *BasicLit) End() token.Pos   { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *ParenExpr) End() token.Pos  { return x.Rparen + 1 }
func (x *IndexExpr) End() token.Pos  { return x.Rbrack + 1 }
func (x *CallExpr) End() token.Pos   { return x.Rparen + 1 }
func (x *UnaryExpr) End() token.Pos  { return x.X.End() }
func (x *BinaryExpr) End() token.Pos { return x.Y.End() }

// exprNode() ensures that only expression/type nodes can be
// assigned to an ExprNode.
//
func (*BadExpr) exprNode()    {}
func (*Ident) exprNode()      {}
func (*BasicLit) exprNode()   {}
func (*ParenExpr) exprNode()  {}
func (*IndexExpr) exprNode()  {}
func (*CallExpr) exprNode()   {}
func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}

var noPos token.Pos

// NewIdent creates a new Ident without position.
// Useful for ASTs generated by code other than the Go parser.
//
func NewIdent(name string) *Ident { return &Ident{noPos, name} }
