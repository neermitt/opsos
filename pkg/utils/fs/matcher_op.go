package fs

func Not(m Matcher) Matcher {
	return &notMatcher{m: m}
}

type notMatcher struct {
	m Matcher
}

func (nm *notMatcher) Match(s string) bool {
	return !nm.m.Match(s)
}

func Or(ms ...Matcher) Matcher {
	return &orMatcher{ms: ms}
}

type orMatcher struct {
	ms []Matcher
}

func (om *orMatcher) Match(s string) bool {
	for _, m := range om.ms {
		if m.Match(s) {
			return true
		}
	}
	return false
}

func And(ms ...Matcher) Matcher {
	return &andMatcher{ms: ms}
}

type andMatcher struct {
	ms []Matcher
}

func (am *andMatcher) Match(s string) bool {
	for _, m := range am.ms {
		if !m.Match(s) {
			return false
		}
	}
	return true
}
