package dim

import (
	"errors"
)

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

func (g *graph) initVisted() map[int]bool {
	out := make(map[int]bool)
	for i := 0; i < len(g.adjs); i++ {
		out[i] = false
	}
	return out
}

func duplicate(input map[int]bool) map[int]bool {
	out := make(map[int]bool)
	for k, v := range input {
		out[k] = v
	}
	return out
}

func (g *graph) TopologicalSort() ([]int, error) {
	out := []int{}
	visted := g.initVisted()
	for i := 0; i < len(g.adjs); i++ {
		if !visted[i] {
			o, err := g.tops(i, visted, g.initVisted())
			if err != nil {
				return nil, err
			}
			out = append(out, o...)
		}
	}
	return out, nil
}

func (g *graph) tops(i int, visted map[int]bool, path map[int]bool) ([]int, error) {
	visted[i] = true
	path[i] = true
	out := []int{}
	for j := 0; j < len(g.adjs); j++ {
		if g.adjs[i][j] {
			if path[j] {
				return nil, errors.New("Cycle detected")
			}
			if !visted[j] {
				o, err := g.tops(j, visted, duplicate(path))
				if err != nil {
					return nil, err
				}
				out = append(out, o...)
			}
		}
	}
	out = append(out, i)
	return out, nil
}

func (g *graph) Connect(from, to int) {
	g.adjs[from][to] = true
}
