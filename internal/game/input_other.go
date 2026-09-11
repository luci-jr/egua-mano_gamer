//go:build !(js && wasm)

package game

func getVirtualKey(keyName string) bool {
	return false
}

func resetVirtualKey(keyName string) {
}
