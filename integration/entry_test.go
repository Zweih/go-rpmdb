package integration

import (
	"testing"

	"github.com/Zweih/go-rpmdb/pkg"
	"github.com/stretchr/testify/assert"
)

func Test_headerImport(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			// found by fuzzer
			name: "negative il",
			data: []byte{0xe3, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, _ = pkg.HeaderImport(tt.data)
			})
		})
	}
}
