package gotagger

const similarityThreshold = 0.85

type candidate struct {
	components words
	score      float32
}

type candidates []candidate

func (c candidate) containsWord(w *word) bool {
	for _, component := range c.components {
		if w.isSimilar(component) {
			return true
		}
	}

	return false
}

func (cs candidates) Len() int           { return len(cs) }
func (cs candidates) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs candidates) Less(i, j int) bool { return cs[i].score > cs[j].score }

func (cs candidates) calcScores(ws words) {
	for i, c := range cs {
		for _, w := range ws {
			if c.components.includes(w) {
				cs[i].score += w.score
			}
		}
	}
}
