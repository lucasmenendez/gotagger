// Package gotagger extracts top occurrences words or group of words
package gotagger

import "sort"

// GetTags function returns list of tags generating ngrams (from trigrams to
// unigrams) and count occurrences. Receives list of tokens,
// language code and limit of tags. Return list of tags and error.
func GetTags(txt []string, c string, max int) (tags [][]string, e error) {
	var t *tagger
	if t, e = newTagger(txt, c); e != nil {
		return tags, e
	}

	s := t.score()
	sort.Sort(byScore(s))

	var l []tag = s
	if len(s) > max {
		l = l[:max]
	}

	for _, i := range l {
		tags = append(tags, i.components)
	}

	return tags, e
}
