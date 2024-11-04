// +build linux

package xid

import "os"

func readPlatformMachineID() (string, error) {
	b, err := os.ReadFile("/etc/machine-id")
	if err != nil || len(b) == 0 {
		b, err = os.ReadFile("/sys/class/dmi/id/product_uuid")
	}
	return string(b), err
}
