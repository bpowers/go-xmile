// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmile_test

import (
	"bufio"
	"bytes"
	"fmt"
	xmile "github.com/bpowers/go-xmile/compat"
	"io/ioutil"
	"os"
	"testing"
	"text/template"
)

const dotTmpl = `
digraph model {
{{range $.Clouds}}
{{.}} [shape=none,label=""]{{end}}

{{range $.Stocks}}
{{.Name}} [shape=box]{{end}}

{{range $.Flows}}
{{.Name}} [shape=circle]{{end}}

{{range $.Auxs}}
{{.Name}} [shape=circle]{{end}}

{{range $.Stocks}}{{with $s := .}}{{range $.OutsC}}
{{$s.Name}} -> {{.}}{{end}}{{end}}{{end}}

{{range $.Auxs}}{{with $s := .}}{{range $.OutsC}}
{{$s.Name}} -> {{.}}{{end}}{{end}}{{end}}
}
`

type DotData struct {
	Clouds []string
	Stocks []VInfo
	Flows  []VInfo
	Auxs   []VInfo
}

// VInfo stores information about a var for use by dot.
//
// stock: only connectors out.
// flow: all 3, connector outs & flow ins and outs.
// aux: only outsc.
type VInfo struct {
	Name  string
	OutsC []string
	OutsF []string
	Ins   []string // only flows have ins
}

func writeDot(f *xmile.File) error {
	w := bufio.NewWriter(os.Stderr)
	defer w.Flush()

	var data DotData

	var buf bytes.Buffer
	tmpl := template.New("model.dot")
	if _, err := tmpl.Parse(dotTmpl); err != nil {
		return fmt.Errorf("Parse(dotTmpl): %s", err)
	}
	if err := tmpl.Execute(&buf, &data); err != nil {
		return fmt.Errorf("Execute: %s", err)
	}

	w.Write(buf.Bytes())
	w.Write([]byte("\n"))

	return nil
}

/*
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
*/

func TestDot(t *testing.T) {
	contents, err := ioutil.ReadFile("../models/pred_prey.stmx")
	if err != nil {
		t.Fatalf("ioutil.ReadFile: %s", err)
	}

	f, err := xmile.ReadFile(contents)
	if err != nil {
		t.Fatalf("xmile.ReadFile: %s", err)
	}

	f.Models[0].Interface = xmile.View{}
	if err := writeDot(f); err != nil {
		t.Fatalf("writeDot: %s", err)
	}
}
