// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmile_test

import (
	"bufio"
	"bytes"
	"fmt"
	xmile "github.com/bpowers/go-xmile/compat"
	"github.com/bpowers/go-xmile/smile"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

func normalizeName(n string) string {
	n = strings.Replace(n, `\n`, "_", -1)
	n = strings.ToLower(n)
	return n
}

func normalizeNames(f *xmile.File) {
	for _, m := range f.Models {
		for _, v := range m.Variables {
			v.Name = normalizeName(v.Name)
		}
	}
}

func varMap(m *xmile.Model) map[string]*xmile.Variable {
	vm := make(map[string]*xmile.Variable)
	for _, v := range m.Variables {
		vm[normalizeName(v.Name)] = v
	}
	return vm
}

func refs(v *xmile.Variable) ([]string, error) {
	expr, err := smile.Parse(v.Name, v.Eqn)
	if err != nil {
		return nil, fmt.Errorf("smile.Parse(%s, '%s'): %s", v.Name, v.Eqn, err)
	}
	outs := make([]string, 0)
	var fnNameNext bool
	smile.Inspect(expr, func(n smile.Node) bool {
		if fnNameNext {
			fnNameNext = false
			return true
		}

		switch e := n.(type) {
		case *smile.CallExpr:
			fnNameNext = true
		case *smile.Ident:
			outs = append(outs, normalizeName(e.Name))
		}
		return true
	})
	return outs, nil
}

func writeDot(f *xmile.File) error {
	normalizeNames(f)

	for _, m := range f.Models {
		vm := varMap(m)
		for _, v := range m.Variables {
			outs, err := refs(v)
			log.Printf("var %s refs %v", v.Name, outs)
			if err != nil {
				return fmt.Errorf("refs(%s,'%s'): %s", v.Name, v.Eqn, err)
			}
			_ = outs
		}
		_ = vm

		var data DotData

		var buf bytes.Buffer
		tmpl := template.New("model.dot")
		if _, err := tmpl.Parse(dotTmpl); err != nil {
			return fmt.Errorf("tmpl.Parse(dotTmpl): %s", err)
		}
		if err := tmpl.Execute(&buf, &data); err != nil {
			return fmt.Errorf("tmpl.Execute: %s", err)
		}

		w := bufio.NewWriter(os.Stderr)
		defer w.Flush()
		w.Write(buf.Bytes())
		w.Write([]byte("\n"))
	}

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
