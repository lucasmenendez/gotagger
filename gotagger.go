// Package gotagger extracts top occurrences words or group of words
package gotagger

import "sort"

// GetTags function returns list of tags generating ngrams (from trigrams to unigrams) and count occurrences.
// Recives list of tokens lists, language code and limit of tags. Return list of tags and error.
func GetTags(text [][]string, lang string, limit int) (tags [][]string, err error) {
	var tokens []string
	for _, sentence := range text {
		tokens = append(tokens, sentence...)
	}

	var ngrms [][]string = ngramsRecursive(tokens, 3)

	var t *tagger
	if t, err = newTagger(ngrms, lang); err != nil {
		return tags, err
	}

	scored := t.score()
	sort.Sort(byScore(scored))

	var limited []tag = scored
	if len(scored) > limit {
		limited = limited[:limit]
	}

	for _, i := range limited {
		tags = append(tags, i.components)
	}

	return tags, err
}
