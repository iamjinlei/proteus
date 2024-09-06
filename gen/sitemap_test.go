package gen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSitemap(t *testing.T) {
	m := NewSitemap("unnote.xyz")
	m.Add("index.html")
	m.Add("/foo/this.html")
	m.Add("/bar/that.html")

	xml, err := m.Gen()
	require.NoError(t, err)
	fmt.Printf("%s\n", string(xml))
}
