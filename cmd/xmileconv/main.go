// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// xmileconv converts between vendor-specific XMILE implementations
// and the current TC Draft Spec.  Currently the only vendor-specific
// implementation is isee's... patches welcome.
package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	usageFirstLine = "Usage: %s [OPTION...] FILE"
	usage          = usageFirstLine + `
Convert between vendor-specific and TC Draft Spec XMILE files.

Options:
`
)

var (
	outFmt string
	inFmt  string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&inFmt, "in", "isee",
		"input format [isee,tc]")
	flag.StringVar(&outFmt, "out", "tc",
		"output format [isee,tc]")

	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "error: one and only one argument required.\n%s\n",
			usageFirstLine)
		os.Exit(1)
	}
}

func main() {
}
