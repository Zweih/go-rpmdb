package integration

import (
	"testing"

	"github.com/Zweih/go-rpmdb/pkg"
	"github.com/stretchr/testify/assert"
)

func TestFileDigest(t *testing.T) {
	tests := []struct {
		algorithm pkg.DigestAlgorithm
		expected  string
	}{
		{
			algorithm: pkg.PGPHASHALGO_MD5,
			expected:  "md5",
		},
		{
			algorithm: pkg.PGPHASHALGO_SHA1,
			expected:  "sha1",
		},
		{
			algorithm: pkg.PGPHASHALGO_RIPEMD160,
			expected:  "ripemd160",
		},
		{
			algorithm: 4,
			expected:  "unknown-digest-algorithm",
		},
		{
			algorithm: pkg.PGPHASHALGO_MD2,
			expected:  "md2",
		},
		{
			algorithm: pkg.PGPHASHALGO_TIGER192,
			expected:  "tiger192",
		},
		{
			algorithm: pkg.PGPHASHALGO_HAVAL_5_160,
			expected:  "haval-5-160",
		},
		{
			algorithm: pkg.PGPHASHALGO_SHA256,
			expected:  "sha256",
		},
		{
			algorithm: pkg.PGPHASHALGO_SHA384,
			expected:  "sha384",
		},
		{
			algorithm: pkg.PGPHASHALGO_SHA512,
			expected:  "sha512",
		},
		{
			algorithm: pkg.PGPHASHALGO_SHA224,
			expected:  "sha224",
		},
		{
			algorithm: 12,
			expected:  "unknown-digest-algorithm",
		},
		// assert against known good values
		{
			algorithm: 1,
			expected:  "md5",
		},
		{
			algorithm: 2,
			expected:  "sha1",
		},
		{
			algorithm: 8,
			expected:  "sha256",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			actual := test.algorithm.String()
			assert.Equal(t, test.expected, actual)
		})
	}
}
