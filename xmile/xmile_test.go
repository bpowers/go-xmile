// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmile_test

import (
	"encoding/xml"
	xmile "github.com/bpowers/go-xmile/compat"
	"io/ioutil"
	"os"
	"testing"
)

func cleanIseeDisplay(d *xmile.Display) {
	d.XMLName.Space = ""
	switch d.XMLName.Local {
	case "text_box", "menu_action":
	case "item":
		d.XMLName.Space = "isee"
	case "story", "chapter", "group", "annotation":
		d.XMLName.Space = "isee"
		d.Content = ""
	default:
		d.Content = ""
	}

	for _, c := range d.Children {
		cleanIseeDisplay(c)
	}
}

func TestRead(t *testing.T) {
	contents, err := ioutil.ReadFile("../models/pred_prey.stmx")
	if err != nil {
		t.Fatalf("ioutil.ReadFile: %s", err)
	}
	var f xmile.File
	if err = xml.Unmarshal(contents, &f); err != nil {
		t.Fatalf("xml.Unmarshal: %s", err)
	}

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
			cleanIseeDisplay(v.Display)
		}
		for _, v := range m.Display.Ents {
			cleanIseeDisplay(v)
		}
		for _, v := range m.Interface.Ents {
			cleanIseeDisplay(v)
		}
		m.Interface.SimDelay = nil
	}

	output, err := xml.MarshalIndent(f, "", "    ")
	if err != nil {
		t.Fatalf("xml.MarshalIndent: %s", err)
	}

	os.Stderr.Write([]byte(xmile.XMLDeclaration + "\n"))
	os.Stderr.Write(output)
	os.Stderr.Write([]byte("\n"))
}
