package wildcard

// Result is the final result of a searching
type Result struct {
	matched []rune
	l       int // the total length of the sample
	b       int
	e       int
}

// Matching tells if the searching matchs or not
func (r *Result) Matching() bool {
	return r.e != -1
}

// AllMatching tells if the searching matchs the whole sample or not
func (r *Result) AllMatching() bool {
	return r.l > 0 && r.b == 0 && r.e == r.l
}

// MatchedRunes returns the matched runes or nil if not matching
func (r *Result) MatchedRunes() []rune {
	if r.Matching() {
		return r.matched
	}
	return nil
}

// Interval returns the interval of the result which represents [begin,end)
func (r *Result) Interval() (begin int, end int) {
	return r.b, r.e
}
