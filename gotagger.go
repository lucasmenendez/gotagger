package gotagger

import (
	"io/ioutil"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/lucasmenendez/gotokenizer"
	"log"
	"os"
	"strconv"
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
		var ts []string = gotokenizer.Words(s)
		tk = append(tk, ts...)
	}

	var cleanR *regexp.Regexp = regexp.MustCompile(`¿|\?|!|¡|\.|,|"|'|:|;|<|>|-`)
	for _, r := range tk {
		var w string = cleanR.ReplaceAllString(r, "")
		w = strings.TrimSpace(w)
		t.words = append(t.words, w)
	}
}

func (t *tagger) delStopwords(ws []tag) (candidates []tag, err error) {
	var f string
	if base := os.Getenv("STOPWORDS"); base != "" {
		f = path.Join(base, t.lang)
	} else {
		f = path.Join(stopwords, t.lang)
	}

	var sws []tag
	if b, err := ioutil.ReadFile(f); err != nil {
		return nil, err
	} else {
		var lbr *regexp.Regexp = regexp.MustCompile(`\n`)
		for _, sw := range lbr.Split(string(b), -1) {

			var w tag = tag{components: []string{strings.TrimSpace(sw)}}
			sws = append(sws, w)
		}
	}

	for _, w := range ws {
		var issw bool = false
		for _, sw := range sws {
			if w.contains(sw) {
				issw = true
				break
			}
		}

		if !issw {
			candidates = append(candidates, w)
		}
	}

	return candidates, nil
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

	var err error
	if candidates, err = t.delStopwords(candidates); err != nil {
		log.Println(err)
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

	var err error
	if candidates, err = t.delStopwords(candidates); err != nil {
		log.Println(err)
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
	t.clean()

	var res []tag
	var simple []tag = t.tagWords()
	var double []tag = t.tagTuples()
	if len(simple)+len(double) == 0 {
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

	var max int
	if rmax := os.Getenv("MAX_KEYWORDS"); rmax != "" {
		if max, _ = strconv.Atoi(rmax); err != nil {
			max = maxKeywords
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

		if len(tags) == max {
			break
		}
	}

	return tags, nil
}
