package disjointsets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinimal(t *testing.T) {
	var uf DisjointSets[string] = NewUnionFind[string]()
	halloHandle := uf.MakeSet("Hallo")
	weltHandle := uf.MakeSet("Welt")
	assert.Equal(t, "Hallo", uf.Find(halloHandle))
	assert.Equal(t, "Welt", uf.Find(weltHandle))
	uf.Union(halloHandle, weltHandle)
	assert.Equal(t, uf.Find(halloHandle), uf.Find(weltHandle))
}

func TestUnionAlreadyConnected(t *testing.T) {
	uf := NewUnionFind[string]()
	a := uf.MakeSet("A")
	b := uf.MakeSet("B")
	uf.Union(a, b)
	uf.Union(a, b) // roots are identical — early return on ln 50
	assert.Equal(t, uf.Find(a), uf.Find(b))
}

func TestUnionRankSwap(t *testing.T) {
	// Union(a,b) makes A root with rank 1.
	// Union(c,a): c has rank 0 < A's rank 1, so a and b are swapped,
	// keeping A as root without incrementing the rank.
	uf := NewUnionFind[string]()
	a := uf.MakeSet("A")
	b := uf.MakeSet("B")
	c := uf.MakeSet("C")
	uf.Union(a, b)
	uf.Union(c, a)
	assert.Equal(t, "A", uf.Find(b))
	assert.Equal(t, "A", uf.Find(c))
}

func TestUnionInvalidHandleTypePanic(t *testing.T) {
	uf := NewUnionFind[string]()
	a := uf.MakeSet("A")
	assert.Panics(t, func() { uf.Union(a, nil) })
}

func TestUnionHandleOwnershipPanic(t *testing.T) {
	uf1 := NewUnionFind[string]()
	uf2 := NewUnionFind[string]()
	a := uf1.MakeSet("A")
	b := uf2.MakeSet("B")
	assert.Panics(t, func() { uf1.Union(a, b) })
}
