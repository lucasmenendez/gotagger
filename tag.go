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

// Returns number of intersections between tags. Receives 'tag'. Returns counter.
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

// Checks if 'tag' contain another. If 'strict' is 'true' checks if tags are equal, else only if contains
// common components. Receives 'tag' and boolean 'strict' flag. Returns boolean.
func (t tag) containsTag(i tag, strict bool) bool {
	if strict && len(t.components) != len(i.components) {
		return false
	}

	var count int = t.count(i)
	if strict {
		return count == len(t.components)
	}
	return count > 0
}

// Checks if 'tag' match with pattern. If 'strict' is 'true' checks if each tag components match, else only if some
// component do. Receives string 'pattern' and boolean 'strict' flag. Returns boolean.
func (t tag) containsPattern(p string, strict bool) bool {
	var rgx *regexp.Regexp = regexp.MustCompile(p)

	var count int = 0
	for _, c := range t.components {
		if rgx.MatchString(c) {
			count++
		}
	}

	if strict {
		return count == len(t.components)
	}
	return count > 0
}

// Checks if 'tag' contain string. If 'strict' is 'true' checks if each tag components contain it, else only if some
// component do. Receives string and boolean 'strict' flag. Returns boolean.
func (t tag) containsString(s string, strict bool) bool {
	var count int

	_s := strings.ToLower(s)
	for _, c := range t.components {
		_c := strings.ToLower(c)
		if _c == _s {
			count++
		}
	}

	if strict {
		return count == len(t.components)
	}
	return count > 0
}
