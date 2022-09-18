package stack

import "github.com/neermitt/opsos/pkg/utils/fs"

func matcher(includedPaths []string, excludedPaths []string) fs.Matcher {
	includeMatcher := globMatchers(includedPaths)
	excludedMatcher := globMatchers(excludedPaths)

	return fs.And(includeMatcher, fs.Not(excludedMatcher))
}

func globMatchers(paths []string) fs.Matcher {
	matchers := make([]fs.Matcher, 0)
	for _, p := range paths {
		matchers = append(matchers, fs.Glob(p))
	}
	return fs.Or(matchers...)
}
