package gotagger

import (
	"path"
	"regexp"
	"io/ioutil"

	"github.com/lucasmenendez/gotokenizer"
	"strings"
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
	text string
	score int
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

func (t *tagger) uniques() (uniques []string) {
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

func (t *tagger) prepare() error {
	t.clean()
	if err := t.delStopwords(); err != nil {
		return err
	}

	return nil
}

func (t *tagger) tag() {
	var tags []tag
	var us = t.uniques()
	for _, u := range us {
		var s int = 0
		for _, w := range t.words {
			if u == w {
				s++
			}
		}

		if s >= 0 {
			var tg tag = tag{text: u, score: s}
			tags = append(tags, tg)
		}
	}
}

func Tag(lang, text string) (tags []string, err error) {
	var t *tagger = &tagger{lang: lang, text:text}
	if err = t.prepare(); err != nil {
		return nil, err
	}

	t.tag()
	return t.tags, nil
}