// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// xmileconv converts between vendor-specific XMILE implementations
// and the current TC Draft Spec.  Currently the only vendor-specific
// implementation is isee's... patches welcome.
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/bpowers/go-xmile/compat"
	"io/ioutil"
	"log"
	"os"
)

const (
	usageFirstLine = "Usage: %s [OPTION...] FILE"
	usage          = usageFirstLine + `
Convert between vendor-specific and TC Draft Spec XMILE files.

If file is not specified, attempts to read from stdin.

Options:
`
)

var (
	outFmt string
	inFmt  string

	validFmts = map[string]bool{
		"isee": true,
		"tc":   true,
	}
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&inFmt, "in", "isee",
		"input format [isee,tc]")
	flag.StringVar(&outFmt, "out", "isee",
		"output format [isee,tc]")

	flag.Parse()

	if _, ok := validFmts[inFmt]; !ok {
		fmt.Fprintf(os.Stderr, "error: input format (\"%s\") not recognized.\n%s\n",
			inFmt, usageFirstLine)
		os.Exit(1)
	} else if _, ok := validFmts[outFmt]; !ok {
		fmt.Fprintf(os.Stderr, "error: output format (\"%s\") not recognized.\n%s\n",
			outFmt, usageFirstLine)
		os.Exit(1)
	} else if flag.NArg() != 1 && isTerminal(0) { // isTerminal means we don't have a file piped in
		fmt.Fprintf(os.Stderr, "error: one and only one argument required.\n%s\n",
			usageFirstLine)
		os.Exit(1)
	}
}

func main() {
	var err error
	var contents []byte
	fname := flag.Arg(0)
	if fname == "" {
		fname = "<stdin>"
		contents, err = ioutil.ReadAll(os.Stdin)
	} else {
		contents, err = ioutil.ReadFile(fname)
	}

	// TODO(bp) implement...
	if inFmt != "isee" || outFmt != "isee" {
		log.Fatalf("error: only isee->isee (roundtripping) is supported so far.")
	}

	if err != nil {
		log.Fatalf("ioutil.ReadFile(%s): %s", fname, err)
	}

	f, err := compat.ReadFile(contents)
	if err != nil {
		log.Fatalf("compat.ReadFile: %s", err)
	}

	output, err := xml.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatalf("xml.MarshalIndent: %s", err)
	}

	os.Stderr.Write([]byte(compat.XMLDeclaration + "\n"))
	os.Stderr.Write(output)
	os.Stderr.Write([]byte("\n"))
}
