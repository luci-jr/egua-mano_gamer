//go:build js && wasm

package game

import (
	"syscall/js"
)

func getVirtualKey(keyName string) bool {
	vk := js.Global().Get("_virtualKeys")
	if !vk.IsUndefined() && !vk.IsNull() {
		val := vk.Get(keyName)
		if !val.IsUndefined() && !val.IsNull() {
			return val.Bool()
		}
	}
	return false
}

func resetVirtualKey(keyName string) {
	vk := js.Global().Get("_virtualKeys")
	if !vk.IsUndefined() && !vk.IsNull() {
		vk.Set(keyName, false)
	}
}
