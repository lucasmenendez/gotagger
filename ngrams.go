package gotagger

func ngrams(t [][]string, n int) (res [][]string) {
	if n > 0 {
		for _, s := range t {
			if n >= len(s) {
				res = append(res, s)
				continue
			}

			for i := 0; i < len(s)-n; i++ {
				res = append(res, s[i:i+n])
			}
		}
	}

	return res
}
