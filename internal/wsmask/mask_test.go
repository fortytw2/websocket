package wsmask

import (
	"encoding/binary"
	"math/bits"
	"strconv"
	"testing"

	"github.com/fortytw2/websocket/internal/test/assert"
)

func TestMask(t *testing.T) {
	t.Parallel()

	key := []byte{0xa, 0xb, 0xc, 0xff}
	key32 := binary.LittleEndian.Uint32(key)
	p := []byte{0xa, 0xb, 0xc, 0xf2, 0xc}
	gotKey32 := Mask(key32, p)

	expP := []byte{0, 0, 0, 0x0d, 0x6}
	assert.Equal(t, "p", expP, p)

	expKey32 := bits.RotateLeft32(key32, -8)
	assert.Equal(t, "key32", expKey32, gotKey32)
}

func BenchmarkMask(b *testing.B) {
	sizes := []int{
		2,
		3,
		4,
		8,
		16,
		32,
		128,
		512,
		4096,
		16384,
	}

	fns := []struct {
		name string
		fn   func(b *testing.B, key [4]byte, p []byte)
	}{
		{
			name: "basic",
			fn: func(b *testing.B, key [4]byte, p []byte) {
				for i := 0; i < b.N; i++ {
					basicMask(key, 0, p)
				}
			},
		},

		{
			name: "nhooyr",
			fn: func(b *testing.B, key [4]byte, p []byte) {
				key32 := binary.LittleEndian.Uint32(key[:])
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					Mask(key32, p)
				}
			},
		},
	}

	key := [4]byte{1, 2, 3, 4}

	for _, size := range sizes {
		p := make([]byte, size)

		b.Run(strconv.Itoa(size), func(b *testing.B) {
			for _, fn := range fns {
				b.Run(fn.name, func(b *testing.B) {
					b.SetBytes(int64(size))

					fn.fn(b, key, p)
				})
			}
		})
	}
}

func basicMask(maskKey [4]byte, pos int, b []byte) int {
	for i := range b {
		b[i] ^= maskKey[pos&3]
		pos++
	}
	return pos & 3
}
