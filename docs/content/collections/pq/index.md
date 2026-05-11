---
title: 'Priority Queues'
weight: 1
summary: The `collections/pq` package provides a generic min-priority-queue interface with handle-based access.
---

The `collections/pq` package provides a generic min-priority queue. Elements of type `E` are
stored with an associated priority of ordered type `K` (called the *key*). The central
interface is `PriorityQueue[E, K]`.

## Interface

```go
type PriorityQueue[E any, K cmp.Ordered] interface {
    Insert(element E, key K) base.Handle
    DecreaseKey(handle base.Handle, key K)
    Remove(handle base.Handle)
    ExtractMin() (element E, key K)
    Min() (element E, key K)
    Key(handle base.Handle) K
    IsEmpty() bool
    Size() int
}
```

### `Insert`

```go
Insert(element E, key K) base.Handle
```

Adds `element` with the given `key` and returns a handle to it. Non-unique keys are
accepted. The handle can be passed to `DecreaseKey`, `Remove`, and `Key`.

### `ExtractMin` / `Min`

```go
ExtractMin() (element E, key K)
Min() (element E, key K)
```

`ExtractMin` removes and returns the element with the smallest key.
`Min` returns the same element without removing it.
Both panic if the queue is empty.

### `DecreaseKey`

```go
DecreaseKey(handle base.Handle, key K)
```

Reduces the key of the element identified by `handle`. Panics if the new key is
greater than the current key.

### `Remove`

```go
Remove(handle base.Handle)
```

Removes an arbitrary element by handle in O(log n).

### `Key`

```go
Key(handle base.Handle) K
```

Returns the current key of the element identified by `handle`.

### Example

```go
pq := NewLeftistHeap[string, int]()
hA := pq.Insert("Alice", 10)
      pq.Insert("Bob",   3)
      pq.Insert("Carol", 7)

e, k := pq.Min()           // "Bob", 3  — not removed
fmt.Println(e, k)

pq.DecreaseKey(hA, 1)      // Alice's key drops from 10 to 1

e, k = pq.ExtractMin()     // "Alice", 1
fmt.Println(e, k)
```

## Implementations

### Leftist Trees

`NewLeftistHeap[E, K]()` returns an implementation based on a **leftist heap** — a
heap-ordered binary tree that maintains the *leftist property*: the left subtree always
has a right-spine length at least as long as the right subtree. This keeps the right spine
short (O(log n)), so all operations that follow the right spine are efficient.

| Operation      | Time       |
|----------------|------------|
| `Insert`       | O(log n)   |
| `ExtractMin`   | O(log n)   |
| `DecreaseKey`  | O(log n)   |
| `Remove`       | O(log n)   |
| `Min`          | O(1)       |
| `Merge`        | O(log n)   |

`BuildLeftistHeap[E, K](elements, keys)` constructs a heap from k elements in O(k) time
using pairwise merging, and returns both the heap and a slice of handles.

```go
pq := pq.NewLeftistHeap[string, int]()

pq, handles := pq.BuildLeftistHeap(
    []string{"Alice", "Bob", "Carol"},
    []int{10, 3, 7},
)
```