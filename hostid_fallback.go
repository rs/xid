// +build !darwin,!linux,!freebsd

package xid

import "errors"

func readPlatformMachineID() (string, error) {
	return "", errors.New("not implemented")
}
