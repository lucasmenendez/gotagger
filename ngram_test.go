package gotagger

import "testing"

func TestNgrams(t *testing.T) {
	var n int = 3
	var input []string = []string{"Go", "(", "often", "referred", "to", "as", "golang", ")", "is", "a", "programming", "language"}

	var expected [][]string = [][]string{
		{"Go", "(", "often"},
		{"(", "often", "referred"},
		{"often", "referred", "to"},
		{"referred", "to", "as"},
		{"to", "as", "golang"},
		{"as", "golang", ")"},
		{"golang", ")", "is"},
		{")", "is", "a"},
		{"is", "a", "programming"},
		{"a", "programming", "language"},
	}

	var results [][]string = ngrams(input, n)
	if len(results) != len(expected) {
		t.Errorf("Expect %d, got %d", len(expected), len(results))
		return
	}

	for i := 0; i < len(results); i++ {
		if len(results[i]) != len(expected[i]) {
			t.Errorf("Expect %d, got %d", len(expected[i]), len(results[i]))
			continue
		}

		for j := 0; j < len(results[i]); j++ {
			if results[i][j] != expected[i][j] {
				t.Errorf("Expect %q, got %q", expected[i][j], results[i][j])
			}
		}
	}
}

func TestNgramsRecursive(t *testing.T) {
	var n int = 3
	var input []string = []string{"Go", "is", "a", "programming", "language"}

	var expected [][]string = [][]string{
		{"Go", "is", "a"},
		{"is", "a", "programming"},
		{"a", "programming", "language"},
		{"Go", "is"},
		{"is", "a"},
		{"a", "programming"},
		{"programming", "language"},
		{"Go"},
		{"is"},
		{"a"},
		{"programming"},
		{"language"},
	}

	var results [][]string = ngramsRecursive(input, n)
	if len(results) != len(expected) {
		t.Errorf("Expect %d, got %d", len(expected), len(results))
		return
	}

	for i := 0; i < len(results); i++ {
		if len(results[i]) != len(expected[i]) {
			t.Errorf("Expect %d, got %d", len(expected[i]), len(results[i]))
			return
		}

		for j := 0; j < len(results[i]); j++ {
			if results[i][j] != expected[i][j] {
				t.Errorf("Expect %q, got %q", expected[i][j], results[i][j])
				return
			}
		}
	}
}
