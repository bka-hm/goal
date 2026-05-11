package disjointsets

import (
	"fmt"

	"gitlab.lrz.de/hm/goal/base"
)

type item[E any] struct {
	value  E
	parent int
	rank   int
}

// handle pairs an element index with the UnionFind it belongs to,
// preventing handles from being used with the wrong structure.
type handle[E any] struct {
	index int
	owner *UnionFind[E]
}

type UnionFind[E any] struct {
	elements []*item[E]
}

func NewUnionFind[E any]() *UnionFind[E] {
	return &UnionFind[E]{}
}

func (u *UnionFind[E]) MakeSet(element E) base.Handle {
	index := len(u.elements)
	u.elements = append(u.elements, &item[E]{value: element, parent: index})
	return handle[E]{index: index, owner: u}
}

func (u *UnionFind[E]) validate(h base.Handle) int {
	uh, ok := h.(handle[E])
	if !ok {
		panic(fmt.Sprintf("disjointsets: invalid handle type %T", h))
	}
	if uh.owner != u {
		panic("disjointsets: handle belongs to a different UnionFind")
	}
	return uh.index
}

func (u *UnionFind[E]) Union(first base.Handle, second base.Handle) {
	rootA := u.root(u.validate(first))
	rootB := u.root(u.validate(second))
	if rootA == rootB {
		return
	}
	a, b := u.elements[rootA], u.elements[rootB]
	if a.rank < b.rank {
		a, b = b, a
	}
	b.parent = a.parent
	if a.rank == b.rank {
		a.rank++
	}
}

// Find returns the value of the set representative for the given handle.
// After a Union, both handles return the same (representative's) value.
func (u *UnionFind[E]) Find(h base.Handle) E {
	return u.elements[u.root(u.validate(h))].value
}

// root returns the index of the set representative, compressing the path as it goes.
func (u *UnionFind[E]) root(index int) int {
	if u.elements[index].parent != index {
		u.elements[index].parent = u.root(u.elements[index].parent)
	}
	return u.elements[index].parent
}
