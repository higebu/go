// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd linux netbsd openbsd

package net

import (
	"context"
	"io"
	"os"
	"syscall"
)

func sockaddrToSCTP(sa syscall.Sockaddr) Addr {
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return &SCTPAddr{IP: sa.Addr[0:], Port: sa.Port}
	case *syscall.SockaddrInet6:
		return &SCTPAddr{IP: sa.Addr[0:], Port: sa.Port, Zone: zoneCache.name(int(sa.ZoneId))}
	}
	return nil
}

func (a *SCTPAddr) family() int {
	if a == nil || len(a.IP) <= IPv4len {
		return syscall.AF_INET
	}
	if a.IP.To4() != nil {
		return syscall.AF_INET
	}
	return syscall.AF_INET6
}

func (a *SCTPAddr) sockaddr(family int) (syscall.Sockaddr, error) {
	if a == nil {
		return nil, nil
	}
	return ipToSockaddr(family, a.IP, a.Port, a.Zone)
}

func (a *SCTPAddr) toLocal(net string) sockaddr {
	return &SCTPAddr{loopbackIP(net), a.Port, a.Zone}
}

func (c *SCTPConn) readFrom(r io.Reader) (int64, error) {
	if n, err, handled := sendFile(c.fd, r); handled {
		return n, err
	}
	return genericReadFrom(c, r)
}

func dialSCTP(ctx context.Context, net string, laddr, raddr *SCTPAddr) (*SCTPConn, error) {
	if testHookDialSCTP != nil {
		return testHookDialSCTP(ctx, net, laddr, raddr)
	}
	return doDialSCTP(ctx, net, laddr, raddr)
}

func doDialSCTP(ctx context.Context, net string, laddr, raddr *SCTPAddr) (*SCTPConn, error) {
	fd, err := internetSocket(ctx, net, laddr, raddr, syscall.SOCK_STREAM, syscall.IPPROTO_SCTP, "dial")
	if err != nil {
		return nil, err
	}
	return newSCTPConn(fd), nil
}

func (ln *SCTPListener) ok() bool { return ln != nil && ln.fd != nil }

func (ln *SCTPListener) accept() (*SCTPConn, error) {
	fd, err := ln.fd.accept()
	if err != nil {
		return nil, err
	}
	return newSCTPConn(fd), nil
}

func (ln *SCTPListener) close() error {
	return ln.fd.Close()
}

func (ln *SCTPListener) file() (*os.File, error) {
	f, err := ln.fd.dup()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func listenSCTP(ctx context.Context, network string, laddr *SCTPAddr) (*SCTPListener, error) {
	fd, err := internetSocket(ctx, network, laddr, nil, syscall.SOCK_STREAM, syscall.IPPROTO_SCTP, "listen")
	if err != nil {
		return nil, err
	}
	return &SCTPListener{fd}, nil
}
