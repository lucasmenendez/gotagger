package gotagger

import (
	"sort"
	"testing"
)

func TestPrepare(t *testing.T) {
	if lang, err := loadLanguage("en"); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		var tags []tag = []tag{
			{[]string{"Go", "or", "golang"}, 0},
			{[]string{"or", "golang", "is"}, 0},
			{[]string{"golang", "is", "a"}, 0},
			{[]string{"is", "a", "programing"}, 0},
			{[]string{"a", "programing", "language"}, 0},
			{[]string{"programing", "language", "."}, 0},
		}
		var tgr *tagger = &tagger{lang, tags}
		tgr.prepare()

		var expected []tag = []tag{
			{[]string{"Go", "or", "golang"}, 0},
			{[]string{"programing", "language"}, 0},
		}

		if len(tgr.tags) != len(expected) {
			t.Errorf("Expected %d, got %d", len(expected), len(tgr.tags))
		}

		for i, t1 := range tgr.tags {
			var t2 tag = expected[i]

			if len(t1.components) != len(t2.components) {
				t.Errorf("Expected %d, got %d", len(t2.components), len(t1.components))
				continue
			}

			for j, c1 := range t1.components {
				var c2 string = t2.components[j]
				if c1 != c2 {
					t.Errorf("Expected %s, got %s", c1, c2)
				}
			}
		}
	}
}

func TestNewTagger(t *testing.T) {
	if _, err := newTagger([][]string{}, "en"); err == nil {
		t.Error("Expected no words error, got nil")
	}

	var words [][]string = [][]string{
		{"Go", "or", "golang"},
		{"or", "golang", "is"},
		{"golang", "is", "a"},
		{"is", "a", "programing"},
		{"a", "programing", "language"},
		{"programing", "language", "."},
	}

	if tgr, err := newTagger(words, "en"); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		var tags []tag = tgr.tags
		var expected []tag = []tag{
			{[]string{"Go", "or", "golang"}, 0},
			{[]string{"programing", "language"}, 0},
		}

		for i, t1 := range tags {
			var t2 tag = expected[i]

			if len(t1.components) != len(t2.components) {
				t.Errorf("Expected %d, got %d", len(t2.components), len(t1.components))
				continue
			}

			for j, c1 := range t1.components {
				var c2 string = t2.components[j]
				if c1 != c2 {
					t.Errorf("Expected %s, got %s", c1, c2)
				}
			}
		}
	}
}

func TestScore(t *testing.T) {
	var words [][]string = [][]string{
		{"Go", "or", "golang"},
		{"or", "golang", "is"},
		{"golang", "is", "a"},
		{"is", "a", "programing"},
		{"a", "programing", "language"},
		{"programing", "language", "."},
		{"Go", "or", "golang"},
		{"or", "golang", "is"},
		{"golang", "is", "a"},
		{"is", "a", "programing"},
		{"a", "programing", "language"},
		{"programing", "language", "."},
	}

	if tgr, err := newTagger(words, "en"); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		var temp []tag = tgr.score()
		sort.Sort(byScore(temp))

		var tags []tag = temp[:2]
		var expected []tag = []tag{
			{[]string{"Go", "or", "golang"}, 1},
			{[]string{"programing", "language"}, 1},
		}

		for i, t1 := range tags {
			var t2 tag = expected[i]

			if len(t1.components) != len(t2.components) {
				t.Errorf("Expected %d, got %d", len(t2.components), len(t1.components))
				continue
			}

			if t1.score != t2.score {
				t.Errorf("Expected %d, got %d", t2.score, t1.score)
				continue
			}
		}
	}
}
