// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmile

import (
	"crypto/rand"
	"encoding/xml"
	"fmt"
)

// File represents the entire contents of a XMILE document.
type File struct {
	XMLName xml.Name `xml:"xmile"`

	Version  string   `xml:"version,attr"`
	Level    int      `xml:"level,attr"`
	Header   Header   `xml:"header"`
	SimSpecs SimSpecs `xml:"sim_specs"`
	Models   []*Model `xml:",omitempty"`
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

// SimSpecs defines the time parameters a given model should be
// simulated with, or the defaults for all models defined in a given
// file.
type SimSpecs struct {
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
	Variables []*Variable `xml:"variables>variable,omitempty"`
	Views     []*View     `xml:"views,omitempty>view,omitempty"`
}

// View is a collection of objects representing the visual structure
// of a model, such as a stock and flow diagram, a causal loop
// diagram, or the iThink interface layer.
type View struct {
	XMLName xml.Name `xml:"view"`
	Name string `xml:"name,attr,omitempty"`
}

// Variable is the definition of a model entity.  Some fields, such as
// Inflows and Outflows are only applicable for certain variable
// types.  The type is determined by the tag name and is stored in
// XMLName.Name.
type Variable struct {
	XMLName  xml.Name
	Name     string `xml:"name,attr"`
	Equation string `xml:"eqn"`
	Inflows  []string `xml:"inflow,omitempty"` // empty for non-stocks
	Outflows []string `xml:"outflow,omitempty"` // empty for non-stocks
	Units    string `xml:"units"`
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
