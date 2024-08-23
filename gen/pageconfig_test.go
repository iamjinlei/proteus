package gen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractPageConfig(t *testing.T) {
	cfg, content, err := extractPageConfig([]byte("+++\ntest: config\n+++\ncontent\n"))
	require.NoError(t, err)
	require.Equal(t, 1, len(cfg.m))
	require.Equal(t, "config", cfg.m["test"])
	require.Equal(t, []byte("content\n"), content)

	cfg, content, err = extractPageConfig([]byte("test config\n"))
	require.NoError(t, err)
	require.Equal(t, 0, len(cfg.m))
	require.Equal(t, []byte("test config\n"), content)

	cfg, content, err = extractPageConfig([]byte("+++\ntest config\n"))
	require.NoError(t, err)
	require.Equal(t, 0, len(cfg.m))
	require.Equal(t, []byte("+++\ntest config\n"), content)
}
