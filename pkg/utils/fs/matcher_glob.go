package fs

import "github.com/bmatcuk/doublestar/v4"

func Glob(pattern string) Matcher {
	return &globMatcher{pattern: pattern}
}

type globMatcher struct {
	pattern string
}

func (gm *globMatcher) Match(s string) bool {
	m, err := doublestar.PathMatch(gm.pattern, s)
	if err != nil {
		return false
	}
	return m
}
