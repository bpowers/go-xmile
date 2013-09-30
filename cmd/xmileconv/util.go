// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux

package main

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	// FIXME(bp) implement this for BSD/Darwin to read from stdin.
	return false
}
