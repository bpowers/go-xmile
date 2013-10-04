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
	"log"
	"reflect"
)

// An XML node
type Node interface {
	node()
}

func (*File) node()      {}
func (*IseePrefs) node() {}
func (*Model) node()     {}
func (*Variable) node()  {}

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
	Display   xmile.View  `xml:"display"`
	Interface xmile.View  `xml:"interface"`
}

// Variable is the definition of a model entity.  Some fields, such as
// Inflows and Outflows are only applicable for certain variable
// types.  The type is determined by the tag name and is stored in
// XMLName.Name.
type Variable struct {
	XMLName xml.Name
	xmile.Variable
	Display *xmile.Display `xml:"display"`
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
	case "text_box", "item":
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
	}

	return f, nil
}

func ConvertToIsee(f *xmile.File) (*File, error) {
	return nil, fmt.Errorf("not implemented")
}

func convertFromIseeField(fin reflect.Value, stripVendorTags bool) (fout reflect.Value, err error) {
	vendorField, ok := fin.Interface().(Node)
	if !ok {
		return fin, nil
	}

	var xfin xmile.Node
	xfin, err = ConvertFromIsee(vendorField, stripVendorTags)
	if err != nil {
		err = fmt.Errorf("ConvertFromIsee(%#v): %s", vendorField, err)
	}
	return reflect.ValueOf(xfin).Elem(), nil
}

func convertFromIseeSlice(fin, fout reflect.Value, stripVendorTags bool) error {
	return nil
}

type valProvider func() reflect.Value

// TODO(bp) f is an interface{} so that any tag can be passed, and the
// corresponding TC xmile tag returned.  Currently, only the root File
// tag is supported.
//
// ConvertFromIsee takes an isee tag and converts it to the current TC
// draft XMILE spec.  If stripVendorTags is true, isee-namespaced tags
// and attributes that would otherwise have been passed through will
// be removed.
func ConvertFromIsee(in Node, stripVendorTags bool) (out xmile.Node, err error) {
	switch in.(type) {
	case *File:
		out = new(xmile.File)
	case *Model:
		out = new(xmile.Model)
	case *Variable:
		out = new(xmile.Variable)
	default:
		return nil, fmt.Errorf("value (%#v) not convertable", in)
	}

	vin := reflect.ValueOf(in).Elem()
	vout := reflect.ValueOf(out).Elem()
	nfield := vin.NumField()
	for i := 0; i < nfield; i++ {
		fmt.Printf("\tfield: %s\n", vin.Type().Field(i).Name)
		fin := vin.Field(i)
		foutty, ok := vout.Type().FieldByName(vin.Type().Field(i).Name)
		if !ok {
			log.Printf("field %s not found on TC struct, skipping",
				vin.Type().Field(i).Name)
			continue
		}
		fout := vout.FieldByName(foutty.Name)
		if fin, err = convertFromIseeField(fin, stripVendorTags); err != nil {
			return nil, fmt.Errorf("convertFromVendorTag: %s", err)
		}

		// TODO(bp) model & interface views
		isInd := false
		outVal := fout
		if fout.Kind() == reflect.Ptr {
			isInd = true
			outVal = fout.Elem()
		}

		switch outVal.Kind() {
		case reflect.Slice:
			if fin.Len() == 0 || fin.IsNil() {
				continue
			}
			e0 := fin.Index(0)
			if e0.Kind() == reflect.Ptr {
				e0 = e0.Elem()
			}
			// FIXME(bp) generalize
			if e0.Type() != reflect.TypeOf(Model{}) {
				log.Printf("slice type not model: %s", e0.Type())
				continue
			}
			models := make([]*xmile.Model, fin.Len())
			modelsV := reflect.ValueOf(models)

			for j := 0; j < fin.Len(); j++ {
				m, _ := fin.Index(j).Interface().(*Model)
				var xm xmile.Node
				fmt.Printf("xmodel\n")
				xm, err = ConvertFromIsee(m, stripVendorTags)
				if err != nil {
					return
				}
				xmodel, _ := xm.(*xmile.Model)
				fmt.Printf("xmodel: %#v\n", xmodel)
				models[j] = xmodel
			}
			fout.Set(modelsV)
		default:
			if isInd {
				fout.Set(fin.Addr())
			} else {
				fout.Set(fin)
			}
		}
	}

	return out, nil
}
