package gotagger

// Generates ngrams from list of tokens. Receive list of tokens and n value. Returns list of ngrams.
func ngrams(t []string, n int) (res [][]string) {
	if n > 0 {
		var limit int = len(t) - n + 1
		for i := 0; i < limit; i++ {
			res = append(res, t[i:i+n])
		}
	}

	return res
}

// Call ngrams function recursive, from n to 1. Receive list of tokens and n value. Returns list of ngrams.
func ngramsRecursive(t []string, n int) (res [][]string) {
	if n > 0 {
		for i := n; i >= 1; i-- {
			res = append(res, ngrams(t, i)...)
		}
	}

	return res
}
