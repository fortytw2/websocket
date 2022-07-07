//go:build !js
// +build !js

package websocket

import (
	"bufio"
	"bytes"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/fortytw2/websocket/internal/test/assert"
)

func TestHeader(t *testing.T) {
	t.Parallel()

	t.Run("lengths", func(t *testing.T) {
		t.Parallel()

		lengths := []int{
			124,
			125,
			126,
			127,

			65534,
			65535,
			65536,
			65537,
		}

		for _, n := range lengths {
			n := n
			t.Run(strconv.Itoa(n), func(t *testing.T) {
				t.Parallel()

				testHeader(t, header{
					payloadLength: int64(n),
				})
			})
		}
	})

	t.Run("fuzz", func(t *testing.T) {
		t.Parallel()

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randBool := func() bool {
			return rand.Int()%2 == 0
		}

		for i := 0; i < 10000; i++ {
			h := header{
				fin:    randBool(),
				rsv1:   randBool(),
				rsv2:   randBool(),
				rsv3:   randBool(),
				opcode: opcode(r.Intn(16)),

				masked:        randBool(),
				maskKey:       r.Uint32(),
				payloadLength: r.Int63(),
			}

			testHeader(t, h)
		}
	})
}

func testHeader(t *testing.T, h header) {
	b := &bytes.Buffer{}
	w := bufio.NewWriter(b)
	r := bufio.NewReader(b)

	err := writeFrameHeader(h, w, make([]byte, 8))
	assert.Success(t, err)

	err = w.Flush()
	assert.Success(t, err)

	h2, err := readFrameHeader(r, make([]byte, 8))
	assert.Success(t, err)

	assert.Equal(t, "header fin", h.fin, h2.fin)
	assert.Equal(t, "header rsv1", h.rsv1, h2.rsv1)
	assert.Equal(t, "header rsv2", h.rsv2, h2.rsv2)
	assert.Equal(t, "header rsv3", h.rsv3, h2.rsv3)
	assert.Equal(t, "header opcode", h.opcode, h2.opcode)

	// if masked we need the mask key otherwise ignorei t
	if h.masked {
		assert.Equal(t, "header masked", h.masked, h2.masked)
		assert.Equal(t, "header maskKey", h.maskKey, h2.maskKey)
	} else {
		assert.Equal(t, "header masked", h.masked, h2.masked)
	}

	assert.Equal(t, "read header", h.payloadLength, h2.payloadLength)
}
