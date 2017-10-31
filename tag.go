package gotagger

type tag struct {
	components []string
	score      int
}

type byScore []tag

func (ts byScore) Len() int           { return len(ts) }
func (ts byScore) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }
func (ts byScore) Less(i, j int) bool { return ts[i].score > ts[j].score }

func distance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}
	return 1
}

func (t1 tag) contains(t2 tag) bool {
	for _, c1 := range t1.components {
		for _, c2 := range t2.components {
			if distance(c1, c2) == 0 {
				return true
			}
		}
	}

	return false
}

func (t1 tag) equals(t2 tag) bool {
	for _, c1 := range t1.components {
		var e bool = false
		for _, c2 := range t2.components {
			e = e || c1 == c2
		}

		if !e {
			return false
		}
	}

	return true
}
