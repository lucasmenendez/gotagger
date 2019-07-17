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

	ws, uws, cs := rake(text, lang)
	uws.calcFrequencies(ws)
	uws.calcDegrees(cs)
	uws.calcScores()
	cs.calcScores(uws)

	tags = compose(cs, max)
	return
}

func bigrams(text []string, lang language) (ws, uws words, cs candidates) {
	var re = regexp.MustCompile(SymbolPattern)

	var _ws []string
	for _, item := range text {
		if !re.MatchString(item) && !lang.isStopword(item) {
			_ws = append(_ws, item)

			var w = &word{component: item}
			if !uws.includes(w) {
				uws = append(uws, w)
			}
		}
	}

	for _, bigram := range ngramsRecursive(_ws, 2) {
		var components words
		for _, component := range bigram {
			components = append(components, &word{component: component, frequency: 1})
		}

		cs = append(cs, candidate{components: components})
	}

	return
}

func rake(text []string, lang language) (ws, uws words, cs candidates) {
	var (
		re      = regexp.MustCompile(SymbolPattern)
		current = candidate{}
	)

	for _, item := range text {
		var next = re.MatchString(item) || lang.isStopword(item)

		if next && len(current.components) > 0 {
			cs = append(cs, current)
			current = candidate{}
		} else if !next {
			var w = &word{component: item, frequency: 1}

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
