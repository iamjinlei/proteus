package gen

import (
	"path/filepath"
	"strings"
)

func normalizeDomain(domain string) string {
	if !strings.HasPrefix(domain, "http://") &&
		!strings.HasPrefix(domain, "https://") {
		domain = "https://" + domain
	}

	return domain
}

func normalizeRelPath(relPath string) string {
	relPath = filepath.Clean(relPath)
	for len(relPath) > 0 && relPath[0] == '/' {
		relPath = relPath[1:]
	}

	return relPath
}
