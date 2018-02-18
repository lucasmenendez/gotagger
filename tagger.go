package gotagger

import (
	"errors"
	"regexp"
	"strings"
)

// Struct to define 'tagger' with tag candidates and language definition.
type tagger struct {
	lang language
	tags []tag
}

// newTagger function constructs a new 'tagger', converting words to tags and
// initializing language struct based on code received. Receives words matrix
// and language code. Return 'tagger' pointer and error.
func newTagger(ws [][]string, code string) (t *tagger, e error) {
	var l language
	if l, e = loadLanguage(code); e != nil {
		return &tagger{}, e
	}

	if len(ws) == 0 {
		return t, errors.New("no words, provided")
	}

	var tags []tag
	for _, w := range ws {
		if len(w) > 0 {
			tags = append(tags, tag{w, 0})
		}
	}

	t = &tagger{l, tags}
	t.prepare()

	return t, e
}

// prepare function delete special symbols from each tag from tagger. Then check
// if each tag contains stopwords in boundary components.
func (t *tagger) prepare() {

	var re *regexp.Regexp = regexp.MustCompile(`(\s|"|\.\.\.|\.|,|:|\(|\)|\[|\]|\{|\}|¿|\?|¡|\!|[0-9]+\.[0-9]+)`)

	var r []tag
	for _, i := range t.tags {
		var _c []string

		for _, c := range i.components {
			if re.MatchString(c) {
				var candidate string = re.ReplaceAllString(c, "")
				if candidate != "" {
					c = candidate
				} else {
					continue
				}
			}

			_c = append(_c, c)
		}

		if len(_c) > 0 {
			i.components = _c
			r = append(r, i)
		}
	}

	var p []tag
	for _, i := range r {
		var in bool = false
		if len(i.components) <= 2 {
			for _, stw := range t.lang.stopwords {
				in = in || i.containsString(stw, false)
			}
		} else {
			var (
				_cs string = strings.ToLower(i.components[0])
				_ce string = strings.ToLower(i.components[len(i.components)-1])
			)

			for _, stopword := range t.lang.stopwords {
				_s := strings.ToLower(stopword)
				in = in || _s == _cs || _s == _ce
			}
		}

		if !in {
			p = append(p, i)
		}
	}

	t.tags = p
}

// score functions scores each tag from tagger, counting its occurrences. First
// obtains unique tags, then score all tags and them assigns scores to same tag
// from unique list. Then weights multi word tag with individual scores.
func (t *tagger) score() (s []tag) {
	var us []tag
	for _, i := range t.tags {
		var in bool = false
		for _, u := range us {
			in = in || u.containsTag(i, true)
		}

		if !in {
			us = append(us, i)
		}
	}

	for i, ti := range t.tags {
		for j, tj := range t.tags {
			if j != i && ti.containsTag(tj, true) {
				ti.score += 1
			}
		}

		if ti.score == 0 {
			continue
		}

		for pos, u := range us {
			if u.score == 0 && u.containsTag(ti, true) {
				s = append(s, ti)
				us[pos].score += ti.score
				break
			}
		}
	}

	for i, ti := range s {
		var lti int = len(ti.components)
		for j, tj := range s {
			if lti > len(tj.components) && ti.containsTag(tj, false) {
				s[i].score += s[j].score
				s[j].score -= ti.count(tj)
			}
		}
	}

	return s
}
