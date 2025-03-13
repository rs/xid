//go:build !darwin && !linux && !freebsd && !openbsd && !windows
// +build !darwin,!linux,!freebsd,!openbsd,!windows

package xid

import "errors"

func readPlatformMachineID() (string, error) {
	return "", errors.New("not implemented")
}
