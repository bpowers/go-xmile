// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// compat provides the ability to read and write XMILE files that
// correspond to the implementation in isee systems STELLA and iThink
// version 10 products.
package compat

import (
	"encoding/xml"
	"fmt"
	"github.com/bpowers/go-xmile/xmile"
)

// the standard XML declaration, declared as a constant for easy
// reuse.
const XMLDeclaration = `<?xml version="1.0" encoding="utf-8" ?>`

// File represents the entire contents of a XMILE document as
// implemented by STELLA & iThink version ~10.0.3
type File struct {
	XMLName    xml.Name           `xml:"http://www.systemdynamics.org/XMILE xmile"`
	Version    string             `xml:"version,attr"`
	Level      int                `xml:"level,attr"`
	IseeHack   string             `xml:"xmlns:isee,attr"` // FIXME(bp) workaround for, I think go issue w.r.t. namespaces
	Header     xmile.Header       `xml:"header"`
	SimSpec    xmile.SimSpec      `xml:"sim_specs"`
	Dimensions []*xmile.Dimension `xml:"dimensions>dim,omitempty"`
	ModelUnits xmile.ModelUnits   `xml:"model_units"`
	IseePrefs  IseePrefs          `xml:"prefs"`
	EqnPrefs   xmile.EqnPrefs     `xml:"equation_prefs"`
	Models     []*Model           `xml:"model,omitempty"`
}

// IseePrefs contains preferences used by STELLA and iThink
type IseePrefs struct {
	XMLName                xml.Name
	Layer                  string         `xml:"layer,attr"`
	GridWidth              string         `xml:"grid_width,attr"`
	GridHeight             string         `xml:"grid_height,attr"`
	DivByZeroAlert         bool           `xml:"divide_by_zero_alert,attr"`
	ShowModPrefix          bool           `xml:"show_module_prefix,attr"`
	HideTransparentButtons bool           `xml:"hide_transparent_buttons,attr"`
	Window                 xmile.Window   `xml:"window"`
	Security               xmile.Security `xml:"security"`
	PrintSetup             xmile.Window   `xml:"print_setup"`
}

// Model represents a container for both the computational definition
// of a system dynamics model, as well as the visual representations
// of that model.
type Model struct {
	XMLName   xml.Name    `xml:"model"`
	Name      string      `xml:"name,attr,omitempty"`
	Variables []*Variable `xml:",any,omitempty"`
	Display   View        `xml:"display"`
	Interface View        `xml:"interface"`
}

// View is a collection of objects representing the visual structure
// of a model, such as a stock and flow diagram, a causal loop
// diagram, or the iThink interface layer.
type View struct {
	XMLName         xml.Name
	Name            string           `xml:"name,attr,omitempty"`
	Ents            []*xmile.Display `xml:",any,omitempty"`
	ScrollX         float64          `xml:"scroll_x,attr"`
	ScrollY         float64          `xml:"scroll_y,attr"`
	Zoom            float64          `xml:"zoom,attr"`
	SimDelay        *SimDelay        `xml:"simulation_delay,omitempty"`
	PageWidth       int              `xml:"page_width,attr,omitempty"`
	PageHeight      int              `xml:"page_height,attr,omitempty"`
	PageRows        int              `xml:"page_rows,attr,omitempty"`
	PageCols        int              `xml:"page_cols,attr,omitempty"`
	PageSequence    string           `xml:"page_sequence,attr,omitempty"`
	ReportFlows     string           `xml:"report_flows,attr,omitempty"`
	ShowPages       bool             `xml:"show_pages,attr,omitempty"` // BUG(bp) default (omitted) when true
	ShowValsOnHover bool             `xml:"show_values_on_hover,attr,omitempty"`
	ConverterSize   string           `xml:"converter_size,attr,omitempty"`
}

// TODO(bp) implement
type SimDelay struct {
}

// Variable is the definition of a model entity.  Some fields, such as
// Inflows and Outflows are only applicable for certain variable
// types.  The type is determined by the tag name and is stored in
// XMLName.Name.
type Variable struct {
	XMLName  xml.Name
	Name     string         `xml:"name,attr"`
	Doc      string         `xml:"doc,omitempty"`
	Eqn      string         `xml:"eqn"`
	NonNeg   *xmile.Exister `xml:"non_negative"`
	Inflows  []string       `xml:"inflow,omitempty"`  // empty for non-stocks
	Outflows []string       `xml:"outflow,omitempty"` // empty for non-stocks
	Units    string         `xml:"units,omitempty"`
	GF       *xmile.GF      `xml:"gf"`
	Display  *xmile.Display `xml:"display"`
}

// NewFile returns a new File object of the given XMILE compliance
// level and name, along with a new UUID.  If you have a file on disk
// you are looking to process, please see ReadFile.
func NewFile(level int, name string) *File {
	id, err := xmile.UUIDv4()
	if err != nil {
		// this is pretty frowned upon, but I don't want
		// NewFile's interface to potentially fail, and if
		// rand.Read fails we have bigger issues.
		panic(err)
	}

	f := &File{Version: "1.0", Level: level}
	f.Header = xmile.Header{
		Name:   name,
		UUID:   id,
		Vendor: "XMILE TC",
		Product: xmile.Product{
			Name:    "go-xmile",
			Version: "0.1",
			Lang:    "en",
		},
	}
	return f
}

// there is a slight impedence mismatch between the spec & the go xml
// marshaler.  cleanIseeDisplayTag works to clean up artifacts related
// to this when reading in <display> tags.
func cleanIseeDisplayTag(d *xmile.Display, inNS bool) {
	d.XMLName.Space = ""
	switch d.XMLName.Local {
	case "text_box", "menu_action", "item":
		// only items with valid content
	case "story":
		inNS = true
		fallthrough
	default:
		d.Content = ""
	}

	if inNS && d.XMLName.Local != "text_box" {
		d.XMLName.Space = "isee"
	}

	for _, c := range d.Children {
		cleanIseeDisplayTag(c, inNS)
	}
}

// ReadFile takes a block of xml content that represents a XMILE file,
// as implemented by iThink/STELLA v10.0.3, and returns a File
// structure, or an error.  The hope is that iThink files rountripped
// through this function will remain readable by iThink.  If not,
// please report it.
func ReadFile(contents []byte) (*File, error) {
	f := new(File)
	if err := xml.Unmarshal(contents, f); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal: %s", err)
	}

	// this bit is cleaning up some weird interactions the go
	// reflection-based code has without isee xmlnamespace.

	// BUG(bp) when we read in a tag with a variable tag name, the
	// XMILE namespace gets propagated to that tag.
	f.IseeHack = "http://iseesystems.com/XMILE"
	f.IseePrefs.XMLName.Space = "isee"
	f.IseePrefs.Window.XMLName.Space = "isee"
	f.IseePrefs.Security.XMLName.Space = "isee"
	f.IseePrefs.PrintSetup.XMLName.Space = "isee"
	f.EqnPrefs.XMLName.Space = "isee"
	for _, m := range f.Models {
		m.Display.XMLName.Space = ""
		m.Interface.XMLName.Space = ""
		for _, v := range m.Variables {
			v.XMLName.Space = ""
			cleanIseeDisplayTag(v.Display, false)
		}
		for _, v := range m.Display.Ents {
			cleanIseeDisplayTag(v, false)
		}
		for _, v := range m.Interface.Ents {
			cleanIseeDisplayTag(v, false)
		}
		m.Interface.SimDelay = nil
	}

	return f, nil
}
