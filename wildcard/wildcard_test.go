package wildcard

import (
	"testing"
)

func TestSearch(t *testing.T) {
	for _, suit := range []struct {
		sample  string
		pattern string
		should  bool
	}{
		{"hehelll", "*he*l", true},
		{"ABCABCCC", "ABCAB??C", true},
		// matching
		{"hello", "hello*", true},
		{"hello", "*hello", true},
		{"hello", "?ello", true},
		{"hello", "hell?", true},
		{"hello", "??*o", true},
		{"hello", "???*", true},
		{"hello", "he**?", true},
		{"hello", "he*?", true},
		// different samples meet the same patterns
		{"hello1", "hello*", true},
		{"hello2", "*hello", false},
		{"hello3", "?ello", false},
		{"hello4", "hell?", false},
		{"hello5", "??*o", false},
		{"hello6", "???*", true},
		{"hello7", "he**?", true},
		{"hello8", "he*?", true},
		// complex pattern
		{"hellohello", "hel?o???*o", true},
		{"hellohello", "?el?o???**", true},
		// escaping
		{"he*llo", "he\\*llo", true},
		// greedy
		{"hehello", "*he*", true},
		{"hehelll", "*he*l", true},
	} {
		fa, _ := Compile([]rune(suit.pattern))
		r := fa.Search([]rune(suit.sample), true)
		if suit.should != r.AllMatching() {
			t.Errorf("on sample %s pattern %s results %v", suit.sample, suit.pattern, r.MatchedRunes())
		}
	}

	for _, suit := range []struct {
		sample  string
		pattern string
		begin   int
		end     int
	}{
		// multi-branchign
		{"hehehol", "heh?l", 2, 7},
		// reluctant
		{"hehello", "*he*", 0, 2},
		{"hehello", "*he*?h", 0, 3},
	} {

		fa, _ := Compile([]rune(suit.pattern))
		r := fa.Search([]rune(suit.sample), false)
		b, e := r.Interval()
		if b != suit.begin && e != suit.end {
			t.Errorf("should be [%d,%d) on pattern %s reluctant searching yet got [%d,%d)",
				suit.begin, suit.end, suit.pattern, b, e)
		}
	}

}

func TestChannel(t *testing.T) {
	fa, _ := Compile([]rune("*he*"))
	sample := []rune("helhellowheat")
	results := []Result{
		Result{[]rune{'h', 'e'}, 0, 0, 2},
		Result{[]rune{'h', 'e'}, 0, 3, 5},
		Result{[]rune{'h', 'e'}, 0, 9, 11},
	}

	schan := make(chan rune)
	go func(schan chan rune) {
		defer close(schan)
		for _, r := range sample {
			schan <- r
		}
	}(schan)

	i := 0
	for r := range fa.Channel(schan) {
		ret := results[i]
		if r.b != ret.b || r.e != ret.e ||
			r.matched[0] != ret.matched[0] ||
			r.matched[1] != ret.matched[1] {
			t.Errorf("expecting %v yet got %v", ret, *r)
		}
		i++
	}
}
