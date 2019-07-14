// Package gotagger provides functions to tag a text. The text must be tokenized.
package gotagger

import (
	"errors"
	"regexp"
	"sort"
)

// SymbolPattern stores all symbol exception to omit
const SymbolPattern = `(\s|"|\.\.\.|\.|,|:|;|\(|\)|\[|\]|\{|\}|¿|\?|¡|\!|[0-9]+\.[0-9]+)`

func GetTags(text []string, code string, max int) (tags [][]string, err error) {
	var lang language
	if lang, err = loadLanguage(code); err != nil {
		return
	}

	if len(text) == 0 {
		return tags, errors.New("no words, provided")
	}

	ws, uws, cs := prepare(text, lang)
	uws.calcFrecuencies(ws)
	uws.calcDegrees(cs)
	uws.calcScores()
	cs.calcScores(uws)

	tags = compose(cs, max)
	return
}

func prepare(text []string, lang language) (ws, uws words, cs candidates) {
	var (
		re      = regexp.MustCompile(SymbolPattern)
		current = candidate{}
	)

	for _, item := range text {
		var w = &word{component: item, frecuency: 1}
		var next = re.MatchString(item) || lang.isStopword(item)

		if next && len(current.components) > 0 {
			cs = append(cs, current)
			current = candidate{}
		} else if !next {
			current.components = append(current.components, w)
			ws = append(ws, w)
			if !uws.includes(w) {
				uws = append(uws, w)
			}
		}
	}

	return
}

func compose(cs candidates, max int) (tags [][]string) {
	sort.Sort(cs)
	for _, c := range cs[:max] {
		var tag []string
		for _, w := range c.components {
			tag = append(tag, w.component)
		}

		tags = append(tags, tag)
	}

	return
}
