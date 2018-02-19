// tag provides struct and functions to operate with multi word concepts,
// compare and sort them
package gotagger

import "strings"

const similarityThreshold float32 = 0.55

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

// strCompare function implements "StrikeAMatch" string similarity algorithm
// created by Simon White. You will find more information about algorithm here:
// http://www.catalysoft.com/articles/StrikeAMatch.html
func strCompare(s1, s2 string) float32 {
	var (
		cs1 []string   = strings.Split(s1, "")
		cs2 []string   = strings.Split(s2, "")
		bs1 [][]string = ngrams(cs1, 2)
		bs2 [][]string = ngrams(cs2, 2)
		u   float32    = float32(len(cs1)+len(cs2)) - 2
		c   int
	)

	for _, p1 := range bs1 {
		for _, p2 := range bs2 {
			if p1[0] == p2[0] && p1[1] == p2[1] {
				c++
			}
		}
	}

	return (2.0 * float32(c)) / u
}

// isSimilar function check if 'tag' provided is similar to current tag.
// If both tags have a similarity coefficient greater than threshold, it will
// return true. Receives a 'tag'. Return boolean.
func (t tag) isSimilar(i tag) bool {
	var _ct string = strings.ToLower(strings.Join(t.components, ""))
	var _ci string = strings.ToLower(strings.Join(i.components, ""))

	var coeff float32 = strCompare(_ct, _ci)
	return coeff > similarityThreshold
}

// count function returns number of intersections between tags. Receives 'tag'.
// Returns counter.
func (t tag) count(i tag) (count int) {
	for _, ct := range t.components {
		_ct := strings.ToLower(ct)
		for _, ci := range i.components {
			_ci := strings.ToLower(ci)
			if _ct == _ci {
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

	var u int = t.count(i)
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
