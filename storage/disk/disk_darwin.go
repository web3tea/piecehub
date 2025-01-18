//go:build darwin
// +build darwin

package disk

func getDirectIOFlag() int {
	return 0
}
