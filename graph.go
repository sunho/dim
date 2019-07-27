package dim

import "errors"

type graph struct {
	adjs [][]bool
}

func newGraph(n int) *graph {
	adjs := [][]bool{}
	for i := 0; i < n; i++ {
		adjs = append(adjs, make([]bool, n))
	}
	return &graph{
		adjs: adjs,
	}
}

func (g *graph) TopologicalSort() ([]int, error) {
	out := []int{}
	visted := make(map[int]bool)
	for i := 0; i < len(g.adjs); i++ {
		o, err := g.tops(i, visted, make(map[int]bool))
		if err != nil {
			return nil, err
		}
		out = append(out, o...)
	}
	out2 := make([]int, 0, len(out))
	for i := len(out) - 1; i >= 0; i-- {
		out2 = append(out2, out[i])
	}
	return out2, nil
}

func (g *graph) tops(i int, visted map[int]bool, path map[int]bool) ([]int, error) {
	visted[i] = true
	path[i] = true
	out := []int{i}
	for j := 0; j < len(g.adjs); j++ {
		if i != j && g.adjs[i][j] {
			if path[j] {
				return nil, errors.New("Cycle detected")
			}
			if !visted[j] {
				o, err := g.tops(j, visted, path)
				if err != nil {
					return nil, err
				}
				out = append(out, o...)
			}
		}
	}
	return out, nil
}

func (g *graph) Connect(from, to int) {
	g.adjs[from][to] = true
}
