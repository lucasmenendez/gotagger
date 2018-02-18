package gotagger

import (
	"sort"
	"testing"
)

func TestCount(t *testing.T) {
	var t1 tag = tag{[]string{"programing", "language"}, 0}
	var t2 tag = tag{[]string{"a", "programing", "language"}, 0}

	if expected, result := 2, t2.count(t1); expected != result {
		t.Errorf("Expected 2, got %d", result)
	}
}

func TestContainsTag(t *testing.T) {
	var t1 tag = tag{[]string{"programing", "language"}, 0}
	var t2 tag = tag{[]string{"a", "programing", "language"}, 0}

	if match := t1.containsTag(t2, false); !match {
		t.Error("Expected true, got false")
	}

	if match := t1.containsTag(t2, true); match {
		t.Error("Expected false, got true")
	}
}

func TestContainsString(t *testing.T) {
	var t1 tag = tag{[]string{"programing", "language"}, 0}
	var s1 string = `programing`
	var s2 string = `golang`

	if match := t1.containsString(s1); !match {
		t.Error("Expected true, got false")
	}


	if match := t1.containsString(s2); match {
		t.Error("Expected false, got true")
	}

}

func TestTagsSorting(t *testing.T) {
	var (
		tags     []tag = []tag{{score: 1}, {score: 3}, {score: 2}}
		expected []tag = []tag{{score: 3}, {score: 2}, {score: 1}}
	)

	sort.Sort(byScore(tags))
	for i, tg := range tags {
		if tg.score != expected[i].score {
			t.Errorf("Expected %d, got %d", tg.score, expected[i].score)
		}
	}
}
