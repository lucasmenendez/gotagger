package gotagger

import (
	"errors"
	"regexp"
	"strings"
)

const SYMBOL_PATTERN string =  `(\s|"|\.\.\.|\.|,|:|;|\(|\)|\[|\]|\{|\}|¿|\?|¡|\!|[0-9]+\.[0-9]+)`

// Struct to define 'tagger' with tag candidates and language definition.
type tagger struct {
	lang language
	pattern string
	uniques []tag
	candidates []tag
}

// newTagger function constructs a new 'tagger', converting words to tags and
// initializing language struct based on code received. Receives words matrix
// and language code. Return 'tagger' pointer and error.
func newTagger(ws []string, code string) (t *tagger, e error) {
	var l language
	if l, e = loadLanguage(code); e != nil {
		return &tagger{}, e
	}

	if len(ws) == 0 {
		return t, errors.New("no words, provided")
	}

	var tags []tag
	var ngrms [][]string = ngramsRecursive(ws, 3)
	for _, w := range ngrms {
		if len(w) > 0 {
			tags = append(tags, tag{w, 0})
		}
	}

	t = &tagger{ lang: l, pattern: SYMBOL_PATTERN }
	var c []tag = t.clean(tags)
	var s []tag = t.simplify(c)
	t.prepare(s)

	return t, e
}

func (t *tagger) clean(tags []tag) (r []tag) {
	var re *regexp.Regexp = regexp.MustCompile(t.pattern)

	for _, i := range tags {
		var _c []string

		var min, max int = 0, len(i.components) - 1
		for p, c := range i.components {
			if re.MatchString(c) {
				if p != min && p != max {
					_c = []string{}
					break
				}

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

	return r
}

func (t *tagger) simplify(r []tag) (s []tag) {
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
			s = append(s, i)
		}
	}

	return s
}

// prepare function delete special symbols from each tag from tagger. Then check
// if each tag contains stopwords in boundary components.
func (t *tagger) prepare(p []tag) {
	var uqs, cdt []tag

	for _, i := range p {
		if len(i.components) == 1 {
			cdt = append(cdt, i)
		}

		var in bool = false
		for _, c := range uqs {
			in = in || c.containsTag(i, true)
		}

		if !in {
			uqs = append(uqs, i)
		}
	}

	t.uniques = uqs
	t.candidates = cdt
}

// score functions scores each tag from tagger, counting its occurrences. First
// obtains unique tags, then score all tags and them assigns scores to same tag
// from unique list. Then weights multi word tag with individual scores.
func (t *tagger) score() (s []tag) {
	for i, ti := range t.candidates {
		for j, tj := range t.candidates {
			if j != i && ti.similar(tj) {
				ti.score += 1
			}
		}

		if ti.score == 0 {
			continue
		}

		for pos, u := range t.uniques {
			if u.similar(ti) {
				t.uniques[pos].score += ti.score
				break
			}
		}
	}

	for _, i := range t.uniques {
		var in bool = false
		for _, j := range s {
			in = in || j.containsTag(i, true)
		}

		if !in {
			s = append(s, i)
		}
	}

	return s
}
