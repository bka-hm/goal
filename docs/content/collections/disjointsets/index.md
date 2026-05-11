---
title: 'Disjoint Sets'
weight: 2
summary: The `collections/disjointsets` package provides a generic interface for disjoint-set (union-find) data structures.
---

The `collections/disjointsets` package provides data structures for maintaining a partition of
elements into disjoint sets. The central interface is `DisjointSets[E]`.

## Interface

```go
type DisjointSets[E any] interface {
    MakeSet(element E) base.Handle
    Union(first, second base.Handle)
    Find(handle base.Handle) E
}
```

### `MakeSet`

```go
MakeSet(element E) base.Handle
```

Adds `element` as a new singleton set and returns a handle to it. The handle is used
to identify the element in subsequent `Union` and `Find` calls.

### `Union`

```go
Union(first, second base.Handle)
```

Merges the two sets that contain the elements identified by `first` and `second`.
If both handles already belong to the same set, the call is a no-op.

### `Find`

```go
Find(handle base.Handle) E
```

Returns the representative element of the set containing the given handle.
After a `Union`, both handles return the same representative value.

### Example

```go
uf := disjointsets.NewUnionFind[string]()
a := uf.MakeSet("Alice")
b := uf.MakeSet("Bob")
c := uf.MakeSet("Carol")

uf.Union(a, b)

fmt.Println(uf.Find(a) == uf.Find(b)) // true  — same set
fmt.Println(uf.Find(a) == uf.Find(c)) // false — different sets

uf.Union(b, c)
fmt.Println(uf.Find(a) == uf.Find(c)) // true  — all merged
```

## Implementations

### Union-Find (Disjoint Set Trees)

`NewUnionFind[E]()` returns an implementation based on disjoint set trees with
**union by rank** and **path compression**. Both operations run in amortized
O(α(n)) time, where α is the inverse Ackermann function — effectively constant
for all practical input sizes.

```go
uf := disjointsets.NewUnionFind[string]()
```