package dim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopogical(t *testing.T) {
	// Connect 0, 1 means service 0 uses service 1
	g := newGraph(4)
	g.Connect(0, 1)
	g.Connect(0, 2)
	g.Connect(0, 3)
	g.Connect(2, 1)
	g.Connect(1, 3)
	idx, err := g.TopologicalSort()
	assert.NoError(t, err)
	expected := []int{3, 1, 2, 0}
	assert.EqualValues(t, expected, idx)
}

func TestTopogicalCycle(t *testing.T) {
	g := newGraph(4)
	g.Connect(0, 1)
	g.Connect(0, 2)
	g.Connect(2, 0)
	_, err := g.TopologicalSort()
	assert.Error(t, err)
}
