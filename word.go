package gotagger

import "strings"

type word struct {
	component                string
	frequency, degree, score float32
}

type words []*word

func (w *word) isSimilar(i *word) bool {
	var (
		_cw   = strings.ToLower(w.component)
		_ci   = strings.ToLower(i.component)
		coeff = strDistance(_cw, _ci)
	)

	return coeff > similarityThreshold
}

func (ws words) includes(w *word) bool {
	for _, i := range ws {
		if i.component == w.component {
			return true
		}
	}

	return false
}

func (ws words) calcFrequencies(text words) {
	for i, wi := range ws {
		for j, wj := range text {
			if i == j {
				continue
			} else if wi.isSimilar(wj) {
				ws[i].frequency++
			}
		}
	}
}

func (ws words) calcDegrees(cs candidates) {
	for i, w := range ws {
		for _, c := range cs {
			if c.containsWord(w) {
				ws[i].degree += float32(len(c.components))
			}
		}
	}
}

func (ws words) calcScores() {
	for i, w := range ws {
		ws[i].score = w.degree / w.frequency
	}
}
