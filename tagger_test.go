package gotagger

import (
	"sort"
	"testing"
)

func TestClean(t *testing.T) {
	if lang, err := loadLanguage("en"); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		var tags []tag = []tag{
			{[]string{"Go"}, 0},
			{[]string{"or"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"is"}, 0},
			{[]string{"a"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"language"}, 0},
			{[]string{"."}, 0},
			{[]string{"Go", "or"}, 0},
			{[]string{"or", "golang"}, 0},
			{[]string{"golang", "is"}, 0},
			{[]string{"is", "a"}, 0},
			{[]string{"a", "programing"}, 0},
			{[]string{"programing", "language"}, 0},
			{[]string{"language", "."}, 0},
			{[]string{"Go", "or", "golang"}, 0},
			{[]string{"or", "golang", "is"}, 0},
			{[]string{"golang", "is", "a"}, 0},
			{[]string{"is", "a", "programing"}, 0},
			{[]string{"a", "programing", "language"}, 0},
			{[]string{"programing", "language", "."}, 0},
		}
		var tgr *tagger = &tagger{lang: lang, pattern: SYMBOL_PATTERN}
		var cleaned []tag = tgr.clean(tags)

		var expected []tag = []tag{
			{[]string{"Go"}, 0},
			{[]string{"or"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"is"}, 0},
			{[]string{"a"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"language"}, 0},
			{[]string{"Go", "or"}, 0},
			{[]string{"or", "golang"}, 0},
			{[]string{"golang", "is"}, 0},
			{[]string{"is", "a"}, 0},
			{[]string{"a", "programing"}, 0},
			{[]string{"programing", "language"}, 0},
			{[]string{"language"}, 0},
			{[]string{"Go", "or", "golang"}, 0},
			{[]string{"or", "golang", "is"}, 0},
			{[]string{"golang", "is", "a"}, 0},
			{[]string{"is", "a", "programing"}, 0},
			{[]string{"a", "programing", "language"}, 0},
			{[]string{"programing", "language"}, 0},
		}

		if len(cleaned) != len(expected) {
			t.Errorf("Expected %d, got %d", len(expected), len(cleaned))
		}

		for p, i := range cleaned {
			var e tag = expected[p]

			if !i.containsTag(e, true) {
				t.Errorf("Expected %q, got %q", e.components, i.components)
			}
		}

		tags = []tag{
			{[]string{"."}, 0},
			{[]string{";"}, 0},
			{[]string{",", "( )"}, 0},
			{[]string{":"}, 0},
			{[]string{"[...]"}, 0},
		}

		cleaned = tgr.clean(tags)
		if len(cleaned) > 0 {
			t.Errorf("Expected empty tags list, got len %d", len(cleaned))
		}
	}
}
func TestSimplify(t *testing.T) {
	if lang, err := loadLanguage("en"); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		var cleaned []tag = []tag{
			{[]string{"Go"}, 0},
			{[]string{"or"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"is"}, 0},
			{[]string{"a"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"language"}, 0},
			{[]string{"Go", "or"}, 0},
			{[]string{"or", "golang"}, 0},
			{[]string{"golang", "is"}, 0},
			{[]string{"is", "a"}, 0},
			{[]string{"a", "programing"}, 0},
			{[]string{"programing", "language"}, 0},
			{[]string{"language"}, 0},
			{[]string{"Go", "or", "golang"}, 0},
			{[]string{"or", "golang", "is"}, 0},
			{[]string{"golang", "is", "a"}, 0},
			{[]string{"is", "a", "programing"}, 0},
			{[]string{"a", "programing", "language"}, 0},
			{[]string{"programing", "language"}, 0},
		}

		var tgr *tagger = &tagger{lang: lang, pattern: SYMBOL_PATTERN}
		var simplified []tag = tgr.simplify(cleaned)

		var expected []tag = []tag{
			{[]string{"Go"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"language"}, 0},
			{[]string{"Go"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"programing", "language"}, 0},
			{[]string{"language"}, 0},
			{[]string{"Go", "or", "golang"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"programing", "language"}, 0},
			{[]string{"programing", "language"}, 0},
		}

		if len(simplified) != len(expected) {
			t.Errorf("Expected %d, got %d", len(expected), len(simplified))
		}

		for p, i := range simplified {
			var e tag = expected[p]
			if !i.containsTag(e, true) {
				t.Errorf("Expected %q, got %q", e.components, i.components)
			}
		}
	}
}
func TestPrepare(t *testing.T) {
	if lang, err := loadLanguage("en"); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		var simplifiyed []tag = []tag{
			{[]string{"Go"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"language"}, 0},
			{[]string{"Go"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"programing", "language"}, 0},
			{[]string{"language"}, 0},
			{[]string{"Go", "or", "golang"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"programing", "language"}, 0},
			{[]string{"programing", "language"}, 0},
		}

		var tgr *tagger = &tagger{lang: lang, pattern: SYMBOL_PATTERN}
		tgr.prepare(simplifiyed)

		var uniques []tag = []tag{
			{[]string{"Go"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"language"}, 0},
			{[]string{"programing", "language"}, 0},
			{[]string{"Go", "or", "golang"}, 0},
		}

		var candidates []tag = []tag{
			{[]string{"Go"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"language"}, 0},
			{[]string{"Go"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"language"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
		}

		if len(uniques) != len(tgr.uniques) {
			t.Errorf("Expected %d, got %d", len(uniques), len(tgr.uniques))
		}

		for p, i := range tgr.uniques {
			var e tag = uniques[p]
			if !i.containsTag(e, true) {
				t.Errorf("Expected %q, got %q", e.components, i.components)
			}
		}

		if len(candidates) != len(tgr.candidates) {
			t.Errorf("Expected %d, got %d", len(candidates), len(tgr.candidates))
		}

		for p, i := range tgr.candidates {
			var e tag = candidates[p]
			if !i.containsTag(e, true) {
				t.Errorf("Expected %q, got %q", e.components, i.components)
			}
		}
	}
}

func TestNewTagger(t *testing.T) {
	if _, err := newTagger([]string{}, "en"); err == nil {
		t.Error("Excected error, got nil")
	}

	var ws []string = []string{"Go", "or", "golang", "is", "a", "programing", "language", "."}
	if tgr, err := newTagger(ws, "en"); err != nil {
		t.Errorf("Excected nil, got %s", err.Error())
	} else {
		var uniques []tag = []tag{
			{[]string{"Go"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"language"}, 0},
			{[]string{"programing", "language"}, 0},
			{[]string{"Go", "or", "golang"}, 0},
		}

		var candidates []tag = []tag{
			{[]string{"Go"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"language"}, 0},
			{[]string{"Go"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
			{[]string{"language"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"golang"}, 0},
			{[]string{"programing"}, 0},
		}

		if len(uniques) != len(tgr.uniques) {
			t.Errorf("Expected %d, got %d", len(uniques), len(tgr.uniques))
		}

		for _, i := range tgr.uniques {
			var in bool = false
			for _, e := range uniques {
				in = in || i.containsTag(e, true)
			}

			if !in {
				t.Errorf("Expected %q", i.components)
			}
		}

		if len(candidates) != len(tgr.candidates) {
			t.Errorf("Expected %d, got %d", len(candidates), len(tgr.candidates))
		}

		for _, i := range tgr.candidates {
			var in bool = false
			for _, e := range candidates {
				in = in || i.containsTag(e, true)
			}

			if !in {
				t.Errorf("Expected %q", i.components)
			}
		}
	}
}

func TestScore(t *testing.T) {
	var ws []string = []string{"Go", "or", "golang", "is", "a", "programing", "language", "."}

	if tgr, err := newTagger(ws, "en"); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		var scored []tag = tgr.score()
		sort.Sort(byScore(scored))

		var expected []tag = []tag{
			{[]string{"Go", "or", "golang"}, 20},
			{[]string{"golang"}, 20},
			{[]string{"programing", "language"}, 8},
			{[]string{"programing"}, 6},
			{[]string{"Go"}, 2},
			{[]string{"language"}, 2},
		}

		for i, _s := range scored {
			var _e tag = expected[i]
			if !_s.containsTag(_e, true) {
				t.Errorf("Expected (%d) %q, got (%d)%q", _e.score, _e.components, _s.score, _s.components)
			}
		}
	}
}
