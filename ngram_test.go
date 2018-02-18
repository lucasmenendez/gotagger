package gotagger

import "testing"

func TestNgrams(t *testing.T) {
	var n int = 3
	var input []string = []string{"Go", "(", "often", "referred", "to", "as", "golang", ")", "is", "a", "programming", "language"}

	var expected [][]string = [][]string{
		[]string{"Go", "(", "often"},
		[]string{"(", "often", "referred"},
		[]string{"often", "referred", "to"},
		[]string{"referred", "to", "as"},
		[]string{"to", "as", "golang"},
		[]string{"as", "golang", ")"},
		[]string{"golang", ")", "is"},
		[]string{")", "is", "a"},
		[]string{"is", "a", "programming"},
		[]string{"a", "programming", "language"},
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
		[]string{"Go", "is", "a"},
		[]string{"is", "a", "programming"},
		[]string{"a", "programming", "language"},
		[]string{"Go", "is"},
		[]string{"is", "a"},
		[]string{"a", "programming"},
		[]string{"programming", "language"},
		[]string{"Go"},
		[]string{"is"},
		[]string{"a"},
		[]string{"programming"},
		[]string{"language"},
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
