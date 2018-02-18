package gotagger

// ngrams function generates ngrams from list of tokens. Receive list of tokens
// and n value. Returns list of ngrams.
func ngrams(t []string, n int) (r [][]string) {
	if n > 0 {
		var l int = len(t) - n + 1
		for i := 0; i < l; i++ {
			r = append(r, t[i:i+n])
		}
	}

	return r
}

// ngramsRecursive functin calls ngrams function recursive, from n to 1. Receive
// list of tokens and n value. Returns list of ngrams.
func ngramsRecursive(t []string, n int) (r [][]string) {
	if n > 0 {
		for i := n; i >= 1; i-- {
			r = append(r, ngrams(t, i)...)
		}
	}

	return r
}
