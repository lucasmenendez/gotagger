package gotagger

import (
	"regexp"
	"strings"
)

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

	var r int = t.count(i)
	if strict {
		return r == len(t.components)
	}
	return r > 0
}

// containsPattern function checks if 'tag' match with pattern. If 'strict' is
// 'true' checks if each tag components match, else only if some component do.
// Receives string 'pattern' and boolean 'strict' flag. Returns boolean.
func (t tag) containsPattern(p string, strict bool) bool {
	var re *regexp.Regexp = regexp.MustCompile(p)

	var r int
	for _, c := range t.components {
		if re.MatchString(c) {
			r++
		}
	}

	if strict {
		return r == len(t.components)
	}
	return r > 0
}

// containsString function checks if 'tag' contain string. If 'strict' is 'true'
// checks if each tag components contain it, else only if some component do.
// Receives string and boolean 'strict' flag. Returns boolean.
func (t tag) containsString(s string, strict bool) bool {
	var r int

	_s := strings.ToLower(s)
	for _, c := range t.components {
		_c := strings.ToLower(c)
		if _c == _s {
			r++
		}
	}

	if strict {
		return r == len(t.components)
	}
	return r > 0
}
