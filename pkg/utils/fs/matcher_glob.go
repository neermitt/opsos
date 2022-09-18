package fs

import (
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

func Glob(pattern string) Matcher {
	return &globMatcher{pattern: pattern}
}

type globMatcher struct {
	pattern string
}

func (gm *globMatcher) Match(s string) bool {
	if filepath.IsAbs(s) {
		s = s[1:]
	}
	m, err := doublestar.PathMatch(gm.pattern, s)
	if err != nil {
		return false
	}
	return m
}
