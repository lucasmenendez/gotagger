package gotagger

import (
	"testing"
	"sort"
)

func TestPrepare(t *testing.T) {
	if lang, err := loadLanguage("en"); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		var tags []tag = []tag{
			tag{ []string{ "Go", "or", "golang" }, 0 },
			tag{ []string{ "or", "golang", "is" }, 0 },
			tag{ []string{ "golang", "is", "a" }, 0 },
			tag{ []string{ "is", "a", "programing" }, 0 },
			tag{ []string{ "a", "programing", "language" }, 0 },
			tag{ []string{ "programing", "language", "." }, 0 },
		}
		var tgr *tagger = &tagger{ lang, tags }
		tgr.prepare()

		var expected []tag = []tag{
			tag{ []string{ "Go", "or", "golang" }, 0 },
			tag{ []string{ "programing", "language" }, 0 },
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
		[]string{ "Go", "or", "golang" },
		[]string{ "or", "golang", "is" },
		[]string{ "golang", "is", "a" },
		[]string{ "is", "a", "programing" },
		[]string{ "a", "programing", "language" },
		[]string{ "programing", "language", "." },
	}

	if tgr, err := newTagger(words, "en"); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		var tags []tag = tgr.tags
		var expected []tag = []tag{
			tag{ []string{ "Go", "or", "golang" }, 0 },
			tag{ []string{ "programing", "language" }, 0 },
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
		[]string{ "Go", "or", "golang" },
		[]string{ "or", "golang", "is" },
		[]string{ "golang", "is", "a" },
		[]string{ "is", "a", "programing" },
		[]string{ "a", "programing", "language" },
		[]string{ "programing", "language", "." },
		[]string{ "Go", "or", "golang" },
		[]string{ "or", "golang", "is" },
		[]string{ "golang", "is", "a" },
		[]string{ "is", "a", "programing" },
		[]string{ "a", "programing", "language" },
		[]string{ "programing", "language", "." },
	}

	if tgr, err := newTagger(words, "en"); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		var temp []tag = tgr.score()
		sort.Sort(byScore(temp))

		var tags []tag = temp[:2]
		var expected []tag = []tag{
			tag{ []string{ "Go", "or", "golang" }, 1 },
			tag{ []string{ "programing", "language" }, 1 },
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

