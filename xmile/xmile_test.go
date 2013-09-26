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
	for _, v := range f.Models[0].Variables {
		v.XMLName.Space = ""
	}

	output, err := xml.MarshalIndent(f, "", "    ")
	if err != nil {
		t.Fatalf("xml.MarshalIndent: %s", err)
	}

	os.Stdout.Write([]byte(xmile.XMLDeclaration + "\n"))
	os.Stdout.Write(output)
	os.Stdout.Write([]byte("\n"))
}
