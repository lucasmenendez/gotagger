package gotagger

import (
	"io/ioutil"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/lucasmenendez/gotokenizer"
)

const (
	stopwords   string = "./stopwords"
	maxKeywords int    = 10
)

type tagger struct {
	lang, text string
	words      []string
	tags       []string
}

func (t *tagger) clean() {
	var tk []string
	for _, s := range gotokenizer.Sentences(t.text) {
		tk = append(tk, s...)
	}

	var cleanR *regexp.Regexp = regexp.MustCompile(`¿|\?|!|¡|\.|,|"|'|:|;|<|>|-`)
	for _, r := range tk {
		var w string = cleanR.ReplaceAllString(r, "")
		w = strings.TrimSpace(w)
		t.words = append(t.words, w)
	}
}

func (t *tagger) delStopwords() error {
	var sws []string

	var f string = path.Join(stopwords, t.lang)
	if b, err := ioutil.ReadFile(f); err != nil {
		return err
	} else {
		var lbr *regexp.Regexp = regexp.MustCompile(`\n`)
		for _, sw := range lbr.Split(string(b), -1) {
			sws = append(sws, strings.TrimSpace(sw))
		}
	}

	var words []string
	for _, w := range t.words {
		var lw = strings.ToLower(w)
		var issw bool = false
		for _, sw := range sws {
			if lw == sw {
				issw = true
				break
			}
		}

		if !issw {
			words = append(words, w)
		}
	}

	t.words = words
	return nil
}

func (t *tagger) prepare() error {
	t.clean()
	if err := t.delStopwords(); err != nil {
		return err
	}

	return nil
}

func (t *tagger) tagWords() (tags []tag) {
	var candidates []tag
	for i, w1 := range t.words {
		var t1 tag = tag{components: []string{w1}}
		for j, w2 := range t.words {
			var t2 tag = tag{components: []string{w2}}
			if i != j && t1.equals(t2) {
				t1.score++
			}
		}
		if t1.score > 0 {
			candidates = append(candidates, t1)
		}
	}

	for _, c := range candidates {
		var s bool = false
		for _, t := range tags {
			if c.equals(t) {
				s = true
				break
			}
		}

		if !s {
			tags = append(tags, c)
		}
	}
	return tags
}

func (t *tagger) tagTuples() (tags []tag) {
	var candidates []tag
	for i := 0; i < len(t.words)-1; i++ {
		var t1 tag = tag{components: []string{t.words[i], t.words[i+1]}}

		for j := 0; j < len(t.words)-1; j++ {
			var t2 tag = tag{components: []string{t.words[j], t.words[j+1]}}

			if i != j && t1.equals(t2) {
				t1.score++
			}
		}

		if t1.score > 0 {
			candidates = append(candidates, t1)
		}
	}

	for _, c := range candidates {
		var s bool = false
		for _, t := range tags {
			if c.equals(t) {
				s = true
				break
			}
		}

		if !s {
			tags = append(tags, c)
		}
	}

	return tags
}

func Tag(lang, text string) (tags []string, err error) {
	var t *tagger = &tagger{lang: lang, text: text}
	if err = t.prepare(); err != nil {
		return nil, err
	}

	var res []tag
	var simple []tag = t.tagWords()
	var double []tag = t.tagTuples()
	if len(simple) + len(double) == 0 {
		return []string{}, nil
	}

	if res = double; len(double) == 0 {
		res = simple
	} else if len(simple) != 0 {
		for _, s := range simple {
			var score int = s.score
			for _, d := range double {
				if d.contains(s) {
					score -= d.score
				}
			}

			if score > 0 {
				res = append(res, s)
			}
		}
	}

	var av int
	for _, tg := range res {
		av += tg.score
	}
	av = int(av / len(res))

	sort.Sort(byScore(res))
	for _, tg := range res {
		if tg.score > av {
			var raw string = strings.Join(tg.components, " ")
			tags = append(tags, raw)
		}

		if len(tags) == maxKeywords {
			break
		}
	}

	return tags, nil
}
