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
	"regexp"
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

var whitespaceRegexp = regexp.MustCompile("[ \t\r\n_]+")

// CanonicalName takes the string in and converts literal newlines
// into underscores, and collapes multiple underscores into a single
// underscore.
func CanonicalName(in string) string {
	return whitespaceRegexp.ReplaceAllString(in, "_")
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
			for _, c := range v.Parameters {
				c.XMLName.Space = ""
				c.To = CanonicalName(c.To)
				c.From = CanonicalName(c.From)
			}
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

func convertFromIseeSlice(fin reflect.Value, stripVendorTags bool) (fout reflect.Value, err error) {
	if fin.Len() == 0 || fin.IsNil() {
		return fin, nil
	}
	e0 := fin.Index(0)
	needsConvert := true

	var slice interface{}

	switch e0.Interface().(type) {
	case *Model:
		slice = make([]*xmile.Model, fin.Len())
	case *Variable:
		slice = make([]*xmile.Variable, fin.Len())
	case *xmile.Dimension:
		slice = make([]*xmile.Dimension, fin.Len())
		needsConvert = false
	default:
		log.Printf("slice type not supported: %s", e0.Type())
		return reflect.ValueOf([]interface{}{}), nil
	}

	for i := 0; i < fin.Len(); i++ {
		var xm xmile.Node
		if needsConvert {
			m, _ := fin.Index(i).Interface().(Node)
			xm, err = ConvertFromIsee(m, stripVendorTags)
		} else {
			xm, _ = fin.Index(i).Interface().(xmile.Node)
		}
		if err != nil {
			return
		}
		switch sl := slice.(type) {
		case []*xmile.Model:
			sl[i] = xm.(*xmile.Model)
		case []*xmile.Variable:
			sl[i] = xm.(*xmile.Variable)
		case []*xmile.Dimension:
			sl[i] = xm.(*xmile.Dimension)
		}
	}
	return reflect.ValueOf(slice), nil

}

type valProvider func() reflect.Value

// TODO(bp) f is an interface{} so that any tag can be passed, and the
// corresponding TC xmile tag returned.  Currently, only the root File
// tag is supported.
//
// TODO(bp) implement stripVendorTags
//
// ConvertFromIsee takes an isee tag and converts it to the current TC
// draft XMILE spec.  If stripVendorTags is true, isee-namespaced tags
// and attributes that would otherwise have been passed through will
// be removed.
func ConvertFromIsee(in Node, stripVendorTags bool) (out xmile.Node, err error) {
	switch n := in.(type) {
	case *File:
		out = new(xmile.File)
	case *Model:
		xm := new(xmile.Model)
		xm.Views = &[]*xmile.View{new(xmile.View), new(xmile.View)}
		*(*xm.Views)[0] = n.Display
		*(*xm.Views)[1] = n.Interface
		(*xm.Views)[0].XMLName.Local = "view"
		(*xm.Views)[1].XMLName.Local = "view"
		(*xm.Views)[1].Name = "interface"
		for _, v := range n.Variables {
			nd := new(xmile.Display)
			*nd = *v.Display
			nd.XMLName.Local = v.XMLName.Local
			nd.Name = v.Name
			(*xm.Views)[0].Ents = append((*xm.Views)[0].Ents, nd)
		}
		out = xm
	case *Variable:
		xv := new(xmile.Variable)
		*xv = n.Variable
		xv.XMLName = n.XMLName
		out = xv
		return
	default:
		return nil, fmt.Errorf("value (%#v) not convertable", in)
	}

	vin := reflect.ValueOf(in).Elem()
	vout := reflect.ValueOf(out).Elem()
	nfield := vin.NumField()
	for i := 0; i < nfield; i++ {
		//log.Printf("\tfield: %s\n", vin.Type().Field(i).Name)
		fin := vin.Field(i)
		foutty, ok := vout.Type().FieldByName(vin.Type().Field(i).Name)
		if !ok {
			//log.Printf("field %s not found on TC struct, skipping",
			//	vin.Type().Field(i).Name)
			continue
		}
		fout := vout.FieldByName(foutty.Name)
		if fin, err = convertFromIseeField(fin, stripVendorTags); err != nil {
			return nil, fmt.Errorf("convertFromVendorTag: %s", err)
		}

		isInd := false
		outVal := fout
		if fout.Kind() == reflect.Ptr {
			isInd = true
			outVal = fout.Elem()
		}

		switch outVal.Kind() {
		case reflect.Slice:
			fin, err = convertFromIseeSlice(fin, stripVendorTags)
			if err != nil {
				log.Printf("convertFromIseeSlice: %s", err)
				continue
			}
			if fin.Len() == 0 || fin.IsNil() {
				continue
			}
			fallthrough
		default:
			if isInd {
				fout.Set(fin.Addr())
			} else {
				fout.Set(fin)
			}
		}
	}

	// update the header so that consumers know we now have TC
	// XMILE
	switch f := out.(type) {
	case *xmile.File:
		f.Header.Vendor = "SDLabs"
		f.Header.Product = xmile.Product{"go-xmile", "0.1", ""}
	}

	return out, nil
}
