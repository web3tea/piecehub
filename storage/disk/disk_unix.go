//go:build linux
// +build linux

package disk

import "syscall"

func getDirectIOFlag() int {
    return syscall.O_DIRECT
}
