package gotagger

import "strings"

const maxLevenshtain = 0.15

type tag struct {
	components []string
	score      int
}

type byScore []tag

func (ts byScore) Len() int           { return len(ts) }
func (ts byScore) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }
func (ts byScore) Less(i, j int) bool { return ts[i].score > ts[j].score }

func min(nums ...int) (min int) {
	min = nums[0]
	for _, n := range nums {
		if n < min {
			min = n
		}
	}

	return min
}

func levenshtain(w1, w2 string) (diff int) {
	var l1, l2 int = len(w1), len(w2)
	var costs [][]int = make([][]int, l1+1)
	for i := 0; i < l1+1; i++ {
		costs[i] = make([]int, l2+1)
	}

	for i := 0; i < l1+1; i++ {
		costs[i][0] = i
	}

	for j := 0; j < l2+1; j++ {
		costs[0][j] = j
	}

	for i := 1; i < l1+1; i++ {
		for j := 1; j < l2+1; j++ {
			var o int = 1
			if w1[i-1] == w2[j-1] {
				o = 0
			}

			costs[i][j] = min(
				costs[i][j-1]+1,
				costs[i-1][j]+1,
				costs[i-1][j-1]+o,
			)
		}
	}

	return costs[l1][l2]
}

func distance(s1, s2 string) float64 {
	s1, s2 = strings.ToLower(s1), strings.ToLower(s2)
	return float64(levenshtain(s1, s2)) / (float64(len(s1)+len(s2)) / 2.0)
}

func (t1 tag) contains(t2 tag) bool {
	for _, c1 := range t1.components {
		for _, c2 := range t2.components {
			if distance(c1, c2) < maxLevenshtain {
				return true
			}
		}
	}
	return false
}

func (t1 tag) equals(t2 tag) bool {
	if len(t1.components) != len(t2.components) {
		return false
	}
	for i := 0; i < len(t1.components); i++ {
		if distance(t1.components[i], t2.components[i]) > maxLevenshtain {
			return false
		}
	}

	return true
}

func union(st, dt []tag) []tag {
	var todel []int

	for i := 0; i < len(st)-1; i++ {
		var c []string = []string{st[i].components[0], st[i+1].components[0]}
		var s tag = tag{components: c}
		for _, d := range dt {
			if d.contains(s) {
				todel = append(todel, i, i+1)
			}
		}
	}

	return append(st, dt...)
}
