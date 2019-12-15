package util

// WcMatch checks if the given sample string matchs the given pattern
func WcMatch(sample string, pattern string) bool {
	samp, patt := []rune(sample), []rune(pattern)
	slen, plen := len(samp), len(patt)
	asterisk, pos := -1, 0 // the index of the `*` and the position of rune in the sample

	for i := 0; i < slen; i++ {
		s := samp[i]
		// when the current pattern rune is the `*`, it wildly matches everything
		if pos < plen && patt[pos] == '*' {
			asterisk = pos // remember the asterisk's position
			pos++          // check next
		} else if pos < plen && (patt[pos] == '?' || patt[pos] == s) { // when it is `?` or matchs the current sample rune
			pos++ // check next
		} else if asterisk != -1 {
			// if can't match the pattern rune after the asterisk, just let the pattern rune
			// stay and move back the sample cursor
			i += asterisk - pos + 1
			pos = asterisk + 1
		} else {
			continue
		}
	}
	for ; pos < plen && patt[pos] == '*'; pos++ {
	}
	return pos == plen
}
