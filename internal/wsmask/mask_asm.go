//go:build !appengine && (amd64 || arm64)

package wsmask

import "golang.org/x/sys/cpu"

// Mask applies the WebSocket masking algorithm to b
// with the given key.
// See https://tools.ietf.org/html/rfc6455#section-5.3
func Mask(key uint32, b []byte) uint32 {
	if len(b) > 0 {
		return maskAsm(&b[0], len(b), key)
	}
	return key
}

//lint:ignore U1000 used in asm
var useAVX2 = cpu.X86.HasAVX2

//go:noescape
func maskAsm(b *byte, len int, key uint32) uint32
