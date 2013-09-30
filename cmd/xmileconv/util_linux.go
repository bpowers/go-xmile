// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"syscall"
	"unsafe"
)

/*
	ATTRIBUTION: This is from util.go in
	code.google.com/p/go.crypto/ssh/terminal
*/

const ioctlReadTermios = syscall.TCGETS

// isTerminal returns true if the given file descriptor is a terminal.
func isTerminal(fd int) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}
