// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmile

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
	XMLName    xml.Name     `xml:"http://www.systemdynamics.org/XMILE xmile"`
	Version    string       `xml:"version,attr"`
	Level      int          `xml:"level,attr"`
	Header     Header       `xml:"header"`
	SimSpec    SimSpec      `xml:"sim_specs"`
	Dimensions []*Dimension `xml:"dimensions,omitempty>dim,omitempty"`
	ModelUnits *ModelUnits  `xml:"model_units"`
	EqnPrefs   *EqnPrefs    `xml:"equation_prefs"`
	Models     []*Model     `xml:"model,omitempty"`
}

type EqnPrefs struct {
	XMLName xml.Name
	OrderBy string `xml:"order_by,attr"`
}

type ModelUnits struct {
}

// Point represents a position in a 2D plane
type Point struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
}

// Size represents an area on a 2D plane
type Size struct {
	Width  float64 `xml:"width,attr,omitempty"`
	Height float64 `xml:"height,attr,omitempty"`
}

// Rect is an area with a position.
type Rect struct {
	Point
	Size
}

type Window struct {
	XMLName xml.Name
	Size
	Orientation string `xml:"orientation,attr,omitempty"`
}

// TODO(bp) implement and document
type Security struct {
	XMLName xml.Name
}

// Header contains metadata about a given XMILE File.
type Header struct {
	Smile   *Smile  `xml:"smile"`
	Name    string  `xml:"name"`
	UUID    string  `xml:"uuid"`
	Vendor  string  `xml:"vendor"`
	Product Product `xml:"product"`
}

type Dimension struct {
	XMLName xml.Name `xml:"dim"`
	Name    string   `xml:"name,attr"`
	Size    string   `xml:"size,attr"`
}

// Smile contains information on the features used in this model.
type Smile struct {
	Version       string   `xml:"version,attr,omitempty"`
	UsesArrays    int      `xml:"uses_arrays,omitempty"`
	UsesQueue     *Exister `xml:"uses_queue"`
	UsesConveyer  *Exister `xml:"uses_conveyer"`
	UsesSubmodels *Exister `xml:"uses_submodels"`
}

// Exister is used as a pointer when we want to make sure an empty tag
// exists.
type Exister string

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
	XMLName         xml.Name
	Name            string     `xml:"name,attr,omitempty"`
	Ents            []*Display `xml:",any,omitempty"`
	ScrollX         float64    `xml:"scroll_x,attr"`
	ScrollY         float64    `xml:"scroll_y,attr"`
	Zoom            float64    `xml:"zoom,attr"`
	PageWidth       int        `xml:"page_width,attr,omitempty"`
	PageHeight      int        `xml:"page_height,attr,omitempty"`
	PageRows        int        `xml:"page_rows,attr,omitempty"`
	PageCols        int        `xml:"page_cols,attr,omitempty"`
	PageSequence    string     `xml:"page_sequence,attr,omitempty"`
	ReportFlows     string     `xml:"report_flows,attr,omitempty"`
	ShowPages       bool       `xml:"show_pages,attr,omitempty"`
	ShowValsOnHover bool       `xml:"show_values_on_hover,attr,omitempty"`
	ConverterSize   string     `xml:"converter_size,attr,omitempty"`
}

// Variable is the definition of a model entity.  Some fields, such as
// Inflows and Outflows are only applicable for certain variable
// types.  The type is determined by the tag name and is stored in
// XMLName.Name.
type Variable struct {
	XMLName  xml.Name
	Name     string   `xml:"name,attr"`
	Doc      string   `xml:"doc,omitempty"`
	Eqn      string   `xml:"eqn"`
	NonNeg   *Exister `xml:"non_negative"`      // stock,(uni-)flow
	Inflows  []string `xml:"inflow,omitempty"`  // empty for non-stocks
	Outflows []string `xml:"outflow,omitempty"` // empty for non-stocks
	Units    string   `xml:"units,omitempty"`
	GF       *GF      `xml:"gf"` // nil if one doesn't exist
}

// GF contains the definition of a graphical function associated with
// a variable.
type GF struct {
	XMLName  xml.Name `xml:"gf"`
	Discrete bool     `xml:"discrete,attr"`
	XPoints  string   `xml:"xpts"`
	YPoints  string   `xml:"ypts"`
	XScale   Scale    `xml:"xscale"`
	YScale   Scale    `xml:"yscale"`
}

type Scale struct {
	Min float64 `xml:"min,attr"`
	Max float64 `xml:"max,attr"`
}

// Style contains information about the visual appearance of a display
// entity.
type Style struct {
	Background  string `xml:"background,attr,omitempty"`
	Color       string `xml:"color,attr,omitempty"`
	FontFamily  string `xml:"font-family,attr,omitempty"`
	FontSize    string `xml:"font-size,attr,omitempty"`
	FontStyle   string `xml:"font-style,attr,omitempty"`
	FontWeight  string `xml:"font-weight,attr,omitempty"`
	TextAlign   string `xml:"text-align,attr,omitempty"`
	TextDeco    string `xml:"text-decoration,attr,omitempty"`
	Margin      string `xml:"margin,attr,omitempty"`
	Padding     string `xml:"padding,attr,omitempty"`
	BorderColor string `xml:"border-color,attr,omitempty"`
	BorderStyle string `xml:"border-style,attr,omitempty"`
	BorderWidth string `xml:"border-width,attr,omitempty"`
}

// BUG(bp) from and to are attributes on graph, subelements on
// connector.  The go marshaler can't handle this inconsistancy.
//
// Display represents a visual entity in a view, such as a stock,
// flow, button or graph.
type Display struct {
	XMLName xml.Name
	Rect
	Style
	UID             string     `xml:"uid,attr,omitempty"`                 // BUG(bp) should this be an int?
	Title           string     `xml:"title,attr,omitempty"`               // graph
	Type            string     `xml:"type,attr,omitempty"`                // pad,graph
	Cleared         string     `xml:"cleared,attr,omitempty"`             // pad,graph
	TimePrecision   int        `xml:"time_precision,attr,omitempty"`      // pad,graph
	PadPageVisible  int        `xml:"visible_index,attr,omitempty"`       // pad
	Appearance      string     `xml:"appearance,attr,omitempty"`          // button,text_box
	ShowGrid        bool       `xml:"show_grid,attr,omitempty"`           // graph
	StyleStr        string     `xml:"style,attr,omitempty"`               // button
	LockText        bool       `xml:"lock_text,attr,omitempty"`           // text_box
	Fill            string     `xml:"fill,attr,omitempty"`                // graphics_frame
	Label           string     `xml:"label,attr,omitempty"`               // stock,flow,aux
	LabelSide       string     `xml:"label_side,omitempty"`               // stock,flow,aux
	LabelAngle      string     `xml:"label_angle,omitempty"`              // stock,flow,aux
	From            string     `xml:"from,omitempty"`                     // connector
	To              string     `xml:"to,omitempty"`                       // connector
	IconOf          string     `xml:"icon_of,attr,omitempty"`             // graph-pad
	PenWidth        int        `xml:"pen_width,attr,omitempty"`           // graph/plot
	Precision       int        `xml:"precision,attr,omitempty"`           // graph/plot,lamp
	Units           string     `xml:"units,attr,omitempty"`               // lamp
	SeperatorK      bool       `xml:"thousands_separator,attr,omitempty"` // lamp
	ShowName        bool       `xml:"show_name,attr,omitempty"`           // lamp
	RetainEndingVal bool       `xml:"retain_ending_value,attr,omitempty"` // lamp
	Zones           *[]Zone    `xml:"zones>zone"`                         // lamp
	ScrollX         float64    `xml:"scroll_x,attr,omitempty"`            // chapter (storytelling)
	ScrollY         float64    `xml:"scroll_y,attr,omitempty"`            // chapter (storytelling)
	EntRef          *EntRef    `xml:"entity,omitempty"`                   // graph/plot
	Points          *[]*Point  `xml:"pts>pt"`                             // flow,connector
	NavAction       *NavAction `xml:"link"`                               // button
	MenuAction      string     `xml:"menu_action,omitempty"`              // button
	Image           *Image     `xml:"image"`                              // graphics_frame
	Children        []*Display `xml:",any,omitempty"`                     // button,popup,lamp,container
	Content         string     `xml:",chardata"`                          // text_box
}

type Zone struct {
	Type  string  `xml:"type,attr"`
	Color string  `xml:"color,attr"`
	Min   float64 `xml:"min"`
	Max   float64 `xml:"max"`
}

type EntRef struct {
	Name    string `xml:"name,attr,omitempty"`
	Content string `xml:",chardata"`
}

type NavAction struct {
	Target string `xml:"target,attr"`
	Point
	Link string `xml:",innerxml"`
}

type Image struct {
	XMLName xml.Name `xml:"image"`
	Size
	Data string `xml:",chardata"`
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
