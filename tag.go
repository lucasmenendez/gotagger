package gotagger

import (
	"regexp"
	"strings"
)


type tag struct {
	components []string
	score      int
}

type byScore []tag

func (ts byScore) Len() int           { return len(ts) }
func (ts byScore) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }
func (ts byScore) Less(i, j int) bool { return ts[i].score > ts[j].score }

func (t tag) countTag(i tag) (count int) {
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

func (t tag) containsTag(i tag, strict bool) bool {
	if strict && len(t.components) != len(i.components) {
		return false
	}

	var count int = t.countTag(i)
	if strict {
		return count == len(t.components)
	}
	return count > 0
}

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