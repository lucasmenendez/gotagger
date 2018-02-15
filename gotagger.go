package gotagger

import "sort"

func GetTags(tokens [][]string, lang string, limit int) (tags [][]string, err error) {
	var (
		unigrams [][]string = ngrams(tokens, 1)
		bigrams [][]string = ngrams(tokens, 2)
		trigrams [][]string  = ngrams(tokens, 3)
	)

	var ngrams [][]string = append(unigrams, append(bigrams, trigrams...)...)

	var t *tagger
	if t, err = newTagger(ngrams, "es"); err != nil {
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
