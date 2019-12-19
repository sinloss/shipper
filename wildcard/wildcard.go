package wildcard

import (
	"bufio"
	"errors"
	"io"
)

const (
	// AllIn is a special greed state meaning all the remaining characters would
	// match in greedy mode. As the greed state of the FA would never be negative
	// so that we could use nagative values to indicate some special greed state.
	AllIn int = -6576767378
	// Fork represents the forking point
	Fork int = -70798275
)

// `feat` would change the representation of a `k`
type feat int

const (
	literal feat = iota
	symbol
)

// `k` is the key to its rune's corresponding state
type k struct {
	feat
	rune
}

var (
	// `any` represents any of the characters
	any k = k{symbol, '\x00'}
	// `fork` represents the forking point, it is not used as the key of the
	// state map, insted it is used by the function `next` to reserve the `Fork`
	// state while returning the real next state number
	forked k = k{symbol, '\x01'}
)

// FA is the dfa of `Knuth Morris Pratt` algorithm along with a nfa backtracking
// when dealing with '?' characters
type FA struct {
	// access the dfa matrix from state instead of character so that
	// we are not limited by a fixed alphabet.
	m []map[k]int
	// all the backed up x states are stored in it for all forking points, the `x`
	// states are resume states that are used to calculate the next state in case
	// of mismatch
	backx map[int]int
	fin   int // final state
	greed int // the state used for greedy matching
}

// Compile compiles a given pattern rune array to the corresponding FA
func Compile(patt []rune) (*FA, error) {
	plen := len(patt)
	if plen == 0 {
		return nil, errors.New("pattern is empty")
	}
	// reserve one more slot for the final state
	fa := &FA{m: make([]map[k]int, plen+1), backx: map[int]int{}}
	fa.fin, fa.greed = cmpl(&fa.m, fa.backx, 0, 0, patt)
	return fa, nil
}

// Search searches the given sample for the first occurance of the pattern
func (fa *FA) Search(sample []rune, greedy bool) *Result {
	j, x, fp := 0, -1, -1
	r := &Result{nil, len(sample), -1, -1}
	for i, s := range sample {
		j, x = fa.advance(j, x, s, &fp) // try advance
		if fa.check(i, j, r) {          // check
			if greedy {
				if fa.greed == AllIn {
					r.e = r.l
					goto FIN
				}
				continue
			}
		FIN:
			r.matched = sample[r.b:r.e]
			return r
		}
	}
	return r
}

// Channel accepts the given channel and search it for all the occurences of the
// pattern. As the given sample is not determined, the scanning is reluctantly
// performed
func (fa *FA) Channel(sample <-chan rune) <-chan *Result {
	rchan := make(chan *Result)
	go func(sample <-chan rune, rchan chan *Result) {
		defer close(rchan)

		i := 0
		j, x, fp, r := 0, -1, -1, &Result{nil, len(sample), -1, -1}
		acc := []rune{}

		for s := range sample {
			j, x = fa.advance(j, x, s, &fp) // try advance
			if fa.check(i, j, r) {          // check
				acc = append(acc, s)
				r.matched = acc
				rchan <- r
				j, x, fp, r =
					0, -1, -1, &Result{nil, len(sample), -1, -1} // reset
				acc = nil // clear accumulated runes
			}
			if r.b != -1 {
				acc = append(acc, s) // accumulate
			}
			i++
		}
	}(sample, rchan)
	return rchan
}

// Scan performs the `Channel` action on the given reader
func (fa *FA) Scan(r io.Reader) <-chan *Result {
	schan := make(chan rune)
	br := bufio.NewReader(r)
	go func(br *bufio.Reader, schan chan rune) {
		defer close(schan)

		for {
			r, _, err := br.ReadRune()
			if err != nil {
				return // stop the goroutine
			}
			schan <- r
		}
	}(br, schan)
	return fa.Channel(schan)
}

func (fa *FA) advance(j int, x int, s rune, fp *int) (int, int) {
	// try advance
	var n int
	var sym k
	c := k{literal, s}
	if j == fa.fin {
		n, sym = next(fa.m, fa.greed, c)
	} else {
		n, sym = next(fa.m, j, c)
	}

	// trunk transaction
	if x == -1 {
		if sym == forked {
			// start a forking process
			x, _ = next(fa.m, fa.backx[j], c)
			// keep the forking point
			*fp = j
		}
		return n, x
	}

	// forked transaction
	x, _ = next(fa.m, x, c)
	switch {
	case n >= j:
		// next state is greater than current state meaning that the state has
		// successfully advanced
		return n, x
	case n <= *fp:
		// next state is less than the forking point meaning that the forking
		// transaction is over, so reset forking point `fp` and resume state
		// `x`
		*fp = -1
		return x, -1
	default:
		// the next state is less than the current state yet greater than forking
		// point meaning the resume state should be the next state
		return x, x
	}
}

func (fa *FA) check(i int, j int, r *Result) bool {
	// keep the starting index in mind
	if j == 1 {
		r.b = i
	}

	if j == fa.fin {
		// keep the final index in mind
		r.e = i + 1 // final state
		return true
	}
	return false
}

// next caculates the next state
func next(m []map[k]int, j int, c k) (int, k) {
	if n, ok := m[j][c]; ok {
		if n == Fork {
			return m[j][any], forked
		}
		return n, c
	}
	// the given `c` is not explicitly specified in the state map, try the `any`
	return m[j][any], any
}

// j represents the state
// x represents the restart state
func cmpl(m *[]map[k]int, backx map[int]int, j int, x int, patt []rune) (
	int, int) {
	escaped, startx := false, x
	for c, p := range patt {
		if escaped {
			escaped = false
			goto SKIP // just skip without incrementing the j
		}
		switch p {
		case '\\': // escaping the next char
			x = other(m, j, x, c, k{literal, patt[c+1]})
			escaped = true
		case '*':
			n := asterisk(j, j+1, c, patt)
			if n == -1 {
				return j, AllIn
			}
			// break the current cmpl by starting a new one as the asterisk has
			// been met and the start position should shift
			return cmpl(m, backx, j, x, patt[n:])
		case '?':
			question(m, backx, j, x, startx)
			x, _ = next(*m, x, any)
		default:
			x = other(m, j, x, c, k{literal, p})
		}
		j++
	SKIP:
	}
	return j, x
}

func asterisk(j int, nextj, c int, patt []rune) int {
	n, l := c+1, len(patt)
	for ; n < l && (patt[n] == '?' || patt[n] == '*'); n++ {
		// swallow all the following consecutive '?' and '*'
	}
	if n >= l {
		return -1
	}
	return n
}

func question(m *[]map[k]int, backx map[int]int, j int, x int, startx int) {
	match(m, j, any)
	// only when at the first state or at any state whose underlying character is
	// identical with the previous and the first state's underlying character
	// should the current state `j` be the same as the resume state `x`
	if j == x {
		return
	}
	for c, v := range (*m)[x] {
		if v != startx {
			(*m)[j][c] = Fork
		}
	}
	backx[j] = x
}

func other(m *[]map[k]int, j int, x int, c int, p k) int {
	match(m, j, p)
	if c != 0 {
		mismatch(m, j, x, p)
		x, _ := next(*m, x, p) // next restart state
		return x
	}
	if p != any {
		(*m)[j][any] = j // the first char spins when mismatch
	}
	return j
}

func match(m *[]map[k]int, j int, p k) {
	capable(m, j)
	(*m)[j][p] = j + 1 // set match case
}

func mismatch(m *[]map[k]int, j int, x int, p k) {
	capable(m, j)
	for c, v := range (*m)[x] {
		if c != p {
			(*m)[j][c] = v // copy mismatch cases
		}
	}
}

func capable(m *[]map[k]int, i int) {
	// ensure capacity
	if l := len(*m); i >= l {
		*m = append(*m, make([]map[k]int, i-l+1)...)
	}
	// ensure map inited
	if (*m)[i] == nil {
		(*m)[i] = map[k]int{}
	}
}
