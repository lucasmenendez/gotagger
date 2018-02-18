package gotagger

import (
	"errors"
	"regexp"
	"strings"
)

// Struct to define 'tagger' with tag candidates and language definition
type tagger struct {
	lang language
	tags []tag
}

// newTagger function constructs a new 'tagger', converting words to tags and initializing language struct based
// on code received. Receives words matrix and language code. Return 'tagger' pointer and error.
func newTagger(words [][]string, code string) (t *tagger, err error) {
	var lang language
	if lang, err = loadLanguage(code); err != nil {
		return &tagger{}, err
	}

	if len(words) == 0 {
		return t, errors.New("no words, provided")
	}

	var tags []tag
	for _, w := range words {
		if len(w) > 0 {
			tags = append(tags, tag{w, 0})
		}
	}

	t = &tagger{lang, tags}
	t.prepare()

	return t, err
}

// prepare funtion delete special symbols from each tag from tagger. Then check if each tag contains stopwords
// in boundary components.
func (t *tagger) prepare() {
	var rgx *regexp.Regexp = regexp.MustCompile(`(\s|"|\.\.\.|\.|,|:|\(|\)|\[|\]|\{|\}|¿|\?|¡|\!|[0-9]+\.[0-9]+)`)

	var replaced []tag
	for _, item := range t.tags {
		var temp []string

		for _, component := range item.components {
			if rgx.MatchString(component) {
				var candidate string = rgx.ReplaceAllString(component, "")
				if candidate != "" {
					component = candidate
				} else {
					continue
				}
			}

			temp = append(temp, component)
		}

		if len(temp) > 0 {
			item.components = temp
			replaced = append(replaced, item)
		}
	}

	var processed []tag
	for _, item := range replaced {
		var contained bool = false
		if len(item.components) <= 2 {
			for _, stopword := range t.lang.stopwords {
				contained = contained || item.containsString(stopword, false)
			}
		} else {
			var (
				_cs string = strings.ToLower(item.components[0])
				_ce string = strings.ToLower(item.components[len(item.components)-1])
			)

			for _, stopword := range t.lang.stopwords {
				_s := strings.ToLower(stopword)
				contained = contained || _s == _cs || _s == _ce
			}
		}

		if !contained {
			processed = append(processed, item)
		}
	}

	t.tags = processed
}

// score functions scores each tag from tagger, counting its occurrences. First obtains unique tags, then score all
// tags and them assigns scores to same tag from unique list. Then weights multi word tag with individual scores.
func (t *tagger) score() (scored []tag) {
	var uniques []tag
	for _, item := range t.tags {
		var exist bool = false
		for _, unique := range uniques {
			exist = exist || unique.containsTag(item, true)
		}

		if !exist {
			uniques = append(uniques, item)
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

		for pos, unique := range uniques {
			if unique.score == 0 && unique.containsTag(ti, true) {
				scored = append(scored, ti)
				uniques[pos].score += ti.score
				break
			}
		}
	}

	for i, ti := range scored {
		for j, tj := range scored {
			if len(ti.components) > len(tj.components) && ti.containsTag(tj, false) {
				scored[i].score += scored[j].score
				scored[j].score -= ti.count(tj)
			}
		}
	}

	return scored
}
