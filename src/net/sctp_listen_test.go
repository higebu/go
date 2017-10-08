// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd linux netbsd openbsd

package net

import (
	"runtime"
	"testing"
)

func (ln *SCTPListener) port() string {
	_, port, err := SplitHostPort(ln.Addr().String())
	if err != nil {
		return ""
	}
	return port
}

var sctpListenerTests = []struct {
	network string
	address string
}{
	{"sctp", ""},
	{"sctp", "0.0.0.0"},
	{"sctp", "::ffff:0.0.0.0"},
	{"sctp", "::"},

	{"sctp", "127.0.0.1"},
	{"sctp", "::ffff:127.0.0.1"},
	{"sctp", "::1"},

	{"sctp4", ""},
	{"sctp4", "0.0.0.0"},
	{"sctp4", "::ffff:0.0.0.0"},

	{"sctp4", "127.0.0.1"},
	{"sctp4", "::ffff:127.0.0.1"},

	{"sctp6", ""},
	{"sctp6", "::"},

	{"sctp6", "::1"},
}

// TestSCTPListener tests both single and double listen to a test
// listener with same address family, same listening address and
// same port.
func TestSCTPListener(t *testing.T) {
	switch runtime.GOOS {
	case "android", "darwin", "dragonfly", "plan9", "solaris", "nacl", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	}

	for _, tt := range sctpListenerTests {
		if !testableListenArgs(tt.network, JoinHostPort(tt.address, "0"), "") {
			t.Logf("skipping %s test", tt.network+" "+tt.address)
			continue
		}

		ln1, err := Listen(tt.network, JoinHostPort(tt.address, "0"))
		if err != nil {
			t.Fatal(err)
		}
		if err := checkFirstListener(tt.network, ln1); err != nil {
			ln1.Close()
			t.Fatal(err)
		}
		ln2, err := Listen(tt.network, JoinHostPort(tt.address, ln1.(*SCTPListener).port()))
		if err == nil {
			ln2.Close()
		}
		if err := checkSecondListener(tt.network, tt.address, err); err != nil {
			ln1.Close()
			t.Fatal(err)
		}
		ln1.Close()
	}
}
