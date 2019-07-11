package gotagger

import (
	"strings"
)

const similarityThreshold = 0.75

// Struct to define 'tag' object that contains its components and its score.
type tag struct {
	components []string
	score      int
}

// Type to sort tags lists.
type byScore []tag

func (ts byScore) Len() int           { return len(ts) }
func (ts byScore) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }
func (ts byScore) Less(i, j int) bool { return ts[i].score > ts[j].score }

// isSimilar function check if 'tag' provided is similar to current tag.
// If both tags have a similarity coefficient greater than threshold, it will
// return true. Receives a 'tag'. Return boolean.
func (t tag) isSimilar(i tag) bool {
	var _ct string = strings.ToLower(strings.Join(t.components, ""))
	var _ci string = strings.ToLower(strings.Join(i.components, ""))

	//var factor float64 = float64(len(t.components) + len(i.components)) / 2

	var coeff = strDistance(_ct, _ci)
	//fmt.Println(coeff / factor < similarityThreshold, coeff, factor, t.components, i.components)
	return coeff > similarityThreshold
}

// count function returns number of intersections between tags. Receives 'tag'.
// Returns counter.
func (t tag) count(i tag) (count int) {
	for _, ct := range t.components {
		for _, ci := range i.components {
			if strDistance(ct, ci) > similarityThreshold {
				count++
				break
			}
		}
	}

	return count
}

// containsTag function checks if 'tag' contain another. If 'strict' is 'true'
// checks if tags are equal, else only if contains common components. Receives
// 'tag' and boolean 'strict' flag. Returns boolean.
func (t tag) containsTag(i tag, strict bool) bool {
	if strict && len(t.components) != len(i.components) {
		return false
	}

	var u = t.count(i)
	return u == len(t.components)
}

// containsString function checks if any 'tag' component contains string.
// Receives string and boolean 'strict' flag. Returns boolean.
func (t tag) containsString(s string) bool {
	var r int
	var _cs string = strings.ToLower(s)
	for _, c := range t.components {
		_c := strings.ToLower(c)
		if _c == _cs {
			r++
		}
	}

	return r > 0
}
