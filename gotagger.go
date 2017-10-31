package gotagger

import (
	"path"
	"regexp"
	"io/ioutil"

	"github.com/lucasmenendez/gotokenizer"
	"strings"
	"fmt"
)

const (
	stopwords string = "./stopwords"
)

type tagger struct {
	lang, text string
	words []string
	tags []string
}

type tag struct {
	components []string
	score int
}

func distance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}
	return 1
}

func (t1 tag) contains(t2 tag) bool {
	for _, c1 := range t1.components {
		for _, c2 := range t2.components {
			if distance(c1, c2) == 0 {
				return true
			}
		}
	}

	return false
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
	for _, w := range t.words {
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
	for n := 0; n < len(t.words) - 1; n++ {
		if t.words[n] != t.words[n + 1] {
			var w string = fmt.Sprintf("%s %s", t.words[n], t.words[n+1])

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

func (t *tagger) tagWords() (tags []tag){
	var us []string = t.uniquesWords()
	for _, u := range us {
		var s int = 0
		for _, w := range t.words {
			if u == w {
				s++
			}
		}

		if s > 1 {
			var tg tag = tag{text: u, score: s}
			tags = append(tags, tg)
		}
	}

	return tags
}

func (t *tagger) tagTuples() (tags []tag){
	var us []string = t.uniquesTuples()
	for _, u := range us {
		var s int = 0
		for n := 0; n < len(t.words) - 1; n++ {
			var w string = fmt.Sprintf("%s %s", t.words[n], t.words[n + 1])
			if u == w {
				s += 2
			}
		}

		if s > 2 {
			var tg tag = tag{components: []string{u}, score: s}
			tags = append(tags, tg)
		}
	}

	return tags
}

func Tag(lang, text string) (tags []string, err error) {
	var t *tagger = &tagger{lang: lang, text:text}
	if err = t.prepare(); err != nil {
		return nil, err
	}

	var simple []tag = t.tagWords()
	var double []tag = t.tagTuples()
	var res []tag = double

	for _, st := range simple {
		var i bool = false
		for _, dt := range double {
			if dt.contains(st) {
				i = true
				break
			}
		}

		if !i {
			res = append(res, st)
		}
	}

	var av int
	for _, tg := range res {
		av += tg.score
	}
	av = int(av/len(res))

	for _, tg := range res {
		if tg.score > av {
			var raw string = strings.Join(tg.components, " ")
			t.tags = append(t.tags, raw)
			fmt.Println(tg.score, raw)
		}
	}

	return t.tags, nil
}