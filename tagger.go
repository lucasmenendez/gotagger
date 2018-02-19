package gotagger

import (
	"errors"
	"regexp"
)

// SYMBOL_PATTERN stores all symbol exception to omit
const SYMBOL_PATTERN string = `(\s|"|\.\.\.|\.|,|:|;|\(|\)|\[|\]|\{|\}|¿|\?|¡|\!|[0-9]+\.[0-9]+)`

// Struct to define 'tagger' with tag candidates and language definition.
type tagger struct {
	lang       language
	pattern    string
	uniques    []tag
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

	t = &tagger{lang: l, pattern: SYMBOL_PATTERN}
	var c []tag = t.clean(tags)
	var s []tag = t.simplify(c)
	t.prepare(s)

	return t, e
}

// clean function remove special characters and symbols from any component of
// each 'tag' received. Receives a 'tag' list. Return a cleaned 'tag' list.
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

// simplify function reduce each 'tag' components deleting stopwords according
// to 'tagger' language. Receives a 'tag' list. Returns a simplified 'tag' list.
func (t *tagger) simplify(r []tag) (s []tag) {
	for _, i := range r {
		var _cs []string
		for _, _c := range i.components {
			if is := t.lang.isStopword(_c); !is {
				_cs = append(_cs, _c)
			} else if is && len(_cs) > 0 {
				_cs = append(_cs, _c)
			}
		}

		if len(_cs) > 0 {
			var lim int
			for lim = len(_cs); lim >= 0 && t.lang.isStopword(_cs[lim-1]); {
				lim--
			}

			if lim > 0 {
				var cs []string = _cs[:lim]
				if len(cs) > 0 {
					i.components = cs
					s = append(s, i)
				}
			}
		}
	}
	return s
}

// prepare function split 'tag' list provided into two lists. One  contains
// uniques 'tag', and the other contains 'tag's with only one component to
// count occurrences of each one faster. Receives a 'tag' list. No return.
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

// score functions scores each candidate 'tag' counting its occurrences. Then,
// add that scores to all similar unique 'tag'. Return a scored 'tag' list.
func (t *tagger) score() (s []tag) {
	for i, ti := range t.candidates {
		for j, tj := range t.candidates {
			if j != i && ti.isSimilar(tj) {
				ti.score += 1
			}
		}

		if ti.score == 0 {
			continue
		}

		for pos, u := range t.uniques {
			if u.isSimilar(ti) {
				t.uniques[pos].score += ti.score
			}
		}
	}

	for _, i := range t.uniques {
		var in bool = false
		for _, j := range s {
			in = in || j.containsTag(i, true)
		}

		if !in && i.score > 0 {
			s = append(s, i)
		}
	}

	return s
}
