package integration

import (
	"testing"

	"github.com/Zweih/go-rpmdb/pkg"
	"github.com/stretchr/testify/assert"
)

func TestFormatFileFlags(t *testing.T) {
	tests := []struct {
		flags    pkg.FileFlags
		expected string
	}{
		// empty
		{
			flags:    0,
			expected: "",
		},
		// check that the formatting works relative to the configured bits
		{
			flags:    pkg.FileFlags(pkg.RPMFILE_CONFIG),
			expected: "c",
		},
		{
			flags:    pkg.FileFlags(pkg.RPMFILE_DOC),
			expected: "d",
		},
		{
			flags:    pkg.FileFlags(pkg.RPMFILE_MISSINGOK),
			expected: "m",
		},
		{
			flags:    pkg.FileFlags(pkg.RPMFILE_NOREPLACE),
			expected: "n",
		},
		{
			flags:    pkg.FileFlags(pkg.RPMFILE_SPECFILE),
			expected: "s",
		},
		{
			flags:    pkg.FileFlags(pkg.RPMFILE_GHOST),
			expected: "g",
		},
		{
			flags:    pkg.FileFlags(pkg.RPMFILE_LICENSE),
			expected: "l",
		},
		{
			flags:    pkg.FileFlags(pkg.RPMFILE_README),
			expected: "r",
		},
		{
			flags:    pkg.FileFlags(pkg.RPMFILE_ARTIFACT),
			expected: "a",
		},
		{
			flags:    pkg.FileFlags(pkg.RPMFILE_CONFIG | pkg.RPMFILE_DOC | pkg.RPMFILE_SPECFILE | pkg.RPMFILE_MISSINGOK | pkg.RPMFILE_NOREPLACE | pkg.RPMFILE_GHOST | pkg.RPMFILE_LICENSE | pkg.RPMFILE_README | pkg.RPMFILE_ARTIFACT),
			expected: "dcsmnglra",
		},
		{
			flags:    pkg.FileFlags(pkg.RPMFILE_DOC | pkg.RPMFILE_ARTIFACT),
			expected: "da",
		},
		// check that the formatting matches relative to verified correct values
		// see helpful examples from: rpm  --dbpath=/var/lib/rpm -qa --queryformat '%{FILEFLAGS:fflags}|%{FILEFLAGS}\n'
		{
			flags:    pkg.FileFlags(89),
			expected: "cmng",
		},
		{
			flags:    pkg.FileFlags(16),
			expected: "n",
		},
		{
			flags:    pkg.FileFlags(64),
			expected: "g",
		},
		{
			flags:    pkg.FileFlags(17),
			expected: "cn",
		},
		{
			flags:    pkg.FileFlags(4096),
			expected: "a",
		},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			assert.Equal(t, test.expected, test.flags.String())
		})
	}
}
