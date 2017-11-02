package gotagger

import (
	"io/ioutil"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/lucasmenendez/gotokenizer"
	"fmt"
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

func (t *tagger) uniquesWords() (uniques []string) {
	for _, rw := range t.words {
		var w string = strings.ToLower(rw)

		var i bool = false
		for _, u := range uniques {
			if w == u {
				i = true
				break
			}
		}

		if !i {
			uniques = append(uniques, w)
		}
	}

	return uniques
}

func (t *tagger) uniquesTuples() (uniques [][]string) {
	for n := 0; n < len(t.words)-1; n++ {
		if rc1, rc2 := t.words[n], t.words[n+1]; rc1 != rc2 {
			var c1, c2 string = strings.ToLower(rc1), strings.ToLower(rc2)

			var i bool = false
			for _, u := range uniques {
				var u1, u2 string = u[0], u[1]
				if u1 == c1 && u2 == c2 {
					i = true
					break
				}
			}

			if !i {
				uniques = append(uniques, []string{c1, c2})
			}
		}
	}
	return uniques
}

func (t *tagger) prepare() error {
	t.clean()
	if err := t.delStopwords(); err != nil {
		return err
	}

	return nil
}

func (t *tagger) tagWords() (tags []tag) {
	var us []string = t.uniquesWords()
	for _, u := range us {
		var s int = 0
		for _, w := range t.words {
			if u == w {
				s++
			}
		}

		if s > 1 {
			var tg tag = tag{components: []string{u}, score: s}
			tags = append(tags, tg)
		}
	}

	return tags
}

func (t *tagger) tagTuples() (tags []tag) {
	var us [][]string = t.uniquesTuples()
	for _, u := range us {
		var ut tag = tag{components: u, score: 0}
		for n := 0; n < len(t.words)-1; n++ {
			var dt tag = tag{components: []string{t.words[n], t.words[n+1]}}

			if ut.equals(dt) {
				ut.score += 2
			}
		}

		if ut.score > 4 {
			tags = append(tags, ut)
		}
	}

	return tags
}

func Tag(lang, text string) (tags []string, err error) {
	var t *tagger = &tagger{lang: lang, text: text}
	if err = t.prepare(); err != nil {
		return nil, err
	}

	var simple []tag = t.tagWords()
	var double []tag = t.tagTuples()
	var res []tag = union(simple, double)

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

			fmt.Println(raw, tg.score)
		}

		if len(tags) == maxKeywords {
			break
		}
	}

	return tags, nil
}
