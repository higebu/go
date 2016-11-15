// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package net

import (
	"os"
	"syscall"
)

func setSCTPNoDelay(fd *netFD, noDelay bool) error {
	if err := fd.incref(); err != nil {
		return err
	}
	defer fd.decref()
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(fd.sysfd, syscall.IPPROTO_SCTP, SCTP_NODELAY, boolint(noDelay)))
}
