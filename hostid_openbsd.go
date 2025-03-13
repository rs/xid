// +build openbsd

package xid

import "syscall"

func readPlatformMachineID() (string, error) {
	return syscall.Sysctl("hw.uuid")
}
