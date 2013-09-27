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

	f, err := xmile.ReadFile(contents)
	if err != nil {
		t.Fatalf("xmile.ReadFile: %s", err)
	}

	output, err := xml.MarshalIndent(f, "", "    ")
	if err != nil {
		t.Fatalf("xml.MarshalIndent: %s", err)
	}

	os.Stderr.Write([]byte(xmile.XMLDeclaration + "\n"))
	os.Stderr.Write(output)
	os.Stderr.Write([]byte("\n"))
}
