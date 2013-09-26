// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// compat provides the ability to read and write XMILE files that
// correspond to the implementation in isee systems STELLA and iThink
// version 10 products.
package compat

import (
	"crypto/rand"
	"encoding/xml"
	"fmt"
)

// the standard XML declaration, declared as a constant for easy
// reuse.
const XMLDeclaration = `<?xml version="1.0" encoding="utf-8" ?>`

// File represents the entire contents of a XMILE document.
type File struct {
	XMLName xml.Name `xml:"xmile"`
	Version string   `xml:"version,attr"`
	Level   int      `xml:"level,attr"`
	Header  Header   `xml:"header"`
	SimSpec SimSpec  `xml:"sim_specs"`
	Models  []*Model `xml:",omitempty"`
}

// Header contains metadata about a given XMILE File.
type Header struct {
	Name    string  `xml:"name"`
	UUID    string  `xml:"uuid"`
	Vendor  string  `xml:"vendor"`
	Product Product `xml:"product"`
}

// Product contains information about the software that created this
// XMILE document.
type Product struct {
	Name    string `xml:",chardata"`
	Version string `xml:"version,attr"`
	Lang    string `xml:"lang,attr"`
}

// SimSpec defines the time parameters a given model should be
// simulated with, or the defaults for all models defined in a given
// file.
type SimSpec struct {
	TimeUnits string  `xml:"time_units,attr,omitempty"`
	Start     float64 `xml:"start"`
	Stop      float64 `xml:"stop"`
	DT        float64 `xml:"dt"`
	Method    string  `xml:"method,omitempty"`
}

// Model represents a container for both the computational definition
// of a system dynamics model, as well as the visual representations
// of that model.
type Model struct {
	XMLName   xml.Name    `xml:"model"`
	Name      string      `xml:"name,attr,omitempty"`
	Variables []*Variable `xml:"variables>variable,omitempty"`
	Views     []*View     `xml:"views,omitempty>view,omitempty"`
}

// View is a collection of objects representing the visual structure
// of a model, such as a stock and flow diagram, a causal loop
// diagram, or the iThink interface layer.
type View struct {
	XMLName xml.Name `xml:"view"`
	Name    string   `xml:"name,attr,omitempty"`
}

// Variable is the definition of a model entity.  Some fields, such as
// Inflows and Outflows are only applicable for certain variable
// types.  The type is determined by the tag name and is stored in
// XMLName.Name.
type Variable struct {
	XMLName  xml.Name
	Name     string   `xml:"name,attr"`
	Equation string   `xml:"eqn"`
	Inflows  []string `xml:"inflow,omitempty"`  // empty for non-stocks
	Outflows []string `xml:"outflow,omitempty"` // empty for non-stocks
	Units    string   `xml:"units"`
}

// UUIDv4 returns a version 4 (random) variant of a UUID, or an error
// if it can not.
func UUIDv4() (string, error) {
	const uuidBytes = 16
	b := make([]byte, uuidBytes)

	n, err := rand.Read(b)
	if err != nil {
		return "", err
	} else if n != uuidBytes {
		return "", fmt.Errorf("rand.Read(): short read  of %d (wanted %d)", n, uuidBytes)
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// NewFile returns a new File object of the given XMILE compliance
// level and name, along with a new UUID.
func NewFile(level int, name string) *File {
	id, err := UUIDv4()
	if err != nil {
		// this is pretty frowned upon, but I don't want
		// NewFile's interface to potentially fail, and if
		// rand.Read fails we have bigger issues.
		panic(err)
	}

	f := &File{Version: "1.0", Level: level}
	f.Header = Header{
		Name:   name,
		UUID:   id,
		Vendor: "XMILE TC",
		Product: Product{
			Name:    "go-xmile",
			Version: "0.1",
			Lang:    "en",
		},
	}
	return f
}
