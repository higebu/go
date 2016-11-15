// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package net

import (
	"runtime"
	"syscall"
)

func setSCTPNoDelay(fd *netFD, noDelay bool) error {
	err := fd.pfd.SetsockoptInt(syscall.IPPROTO_SCTP, SCTP_NODELAY, boolint(noDelay))
	runtime.KeepAlive(fd)
	return wrapSyscallError("setsockopt", err)
}
