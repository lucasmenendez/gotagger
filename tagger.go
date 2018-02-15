package gotagger

import (
	"regexp"
	"strings"
)

type tagger struct {
	lang language
	tags []tag
}

func newTagger(words [][]string, code string) (t *tagger, err error) {
	var lang language
	if lang, err = loadLanguage(code); err != nil {
		return &tagger{}, err
	}

	var tags []tag
	for _, w := range words {
		if len(w) > 0 {
			tags = append(tags, tag{ w, 0 })
		}
	}

	t = &tagger{ lang, tags }
	t.clean()

	return t, err
}

func (t *tagger) clean() {
	var rgx *regexp.Regexp = regexp.MustCompile(`("|\.\.\.|\.|,|:|\(|\)|\[|\]|\{|\}|¿|\?|¡|\!|[0-9]+\.[0-9]+)`)

	var replaced []tag
	for _, item := range t.tags {
		var ok bool = true
		for j, component := range item.components {
			if rgx.MatchString(component) {
				temp := rgx.ReplaceAllString(component, "")
				if len(temp) == 0 {
					ok = false
				} else {
					item.components[j] = temp
				}
			}

			if !ok { break }
		}

		if ok {
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
				_ce string = strings.ToLower(item.components[len(item.components) - 1])
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
				scored[j].score -= ti.countTag(tj)
			}
		}
	}

	return scored
}