package pq

import (
	"cmp"

	"gitlab.lrz.de/hm/goal/base"
)

// A LeftistHeap is an implementation of collection.PriorityQueue
type LeftistHeap[E any, K cmp.Ordered] struct {
	min  *leftistNode[E, K]
	size int // must be maintained by public methods
}

type leftistNode[E any, K cmp.Ordered] struct {
	element         E
	key             K
	dist            int // 1 + dist(right child); 1 if right is nil
	left, right, up *leftistNode[E, K]
}

// distOrZero returns the dist (length of rightmost path to nil) of this node with nil-safety.
func (node *leftistNode[E, K]) distOrZero() int {
	if node == nil {
		return 0
	}
	return node.dist
}

func NewLeftistHeap[E any, K cmp.Ordered]() *LeftistHeap[E, K] {
	return &LeftistHeap[E, K]{}
}

func (heap *LeftistHeap[E, K]) IsEmpty() bool {
	return heap.min == nil
}

func (heap *LeftistHeap[E, K]) Size() int {
	return heap.size
}

func (heap *LeftistHeap[E, K]) Key(handle base.Handle) K {
	node, ok := handle.(*leftistNode[E, K])
	if !ok {
		panic("pq: invalid handle type")
	}
	return node.key
}

// Internal operation to merge two leftist trees given as root nodes.
// Either or both arguments may be nil; the result is nil iff both are nil.
// Size is not maintained here. dist values and left/right orientation are
// corrected along the merge path.
//
// The returned root's up pointer is NOT updated — the caller is responsible
// for setting it (e.g. to nil when the result becomes heap.min).
func (node *leftistNode[E, K]) merge(other *leftistNode[E, K]) *leftistNode[E, K] {
	// trivial cases
	if node == nil {
		return other
	}
	if other == nil {
		return node
	}

	// implicit selection of the minimum as "node" by swapping arguments
	if cmp.Compare(node.key, other.key) > 0 {
		return other.merge(node)
	}

	// the result of the recursive merge is new right successor
	node.right = node.right.merge(other)
	if node.right != nil {
		node.right.up = node
	}

	// fix left/right and own dist value
	if node.right.distOrZero() > node.left.distOrZero() {
		node.left, node.right = node.right, node.left
	}
	node.dist = node.right.distOrZero() + 1
	return node
}

// Min returns the minimum element (with key/priority).
func (heap *LeftistHeap[E, K]) Min() (E, K) {
	if heap.IsEmpty() {
		panic("pq: heap is empty")
	}
	return heap.min.element, heap.min.key
}

// Insert adds a new element with a given priority (key).
func (heap *LeftistHeap[E, K]) Insert(element E, key K) base.Handle {
	newNode := &leftistNode[E, K]{element: element, key: key, dist: 1}
	heap.min = heap.min.merge(newNode)
	heap.size++
	return newNode
}

// ExtractMin extracts the minimum element (with key/priority).
func (heap *LeftistHeap[E, K]) ExtractMin() (E, K) {
	if heap.IsEmpty() {
		panic("pq: heap is empty")
	}
	returnValue := heap.min
	heap.min = heap.min.left.merge(heap.min.right)
	if heap.min != nil {
		heap.min.up = nil
	}
	heap.size--
	return returnValue.element, returnValue.key
}

// DecreaseKey reduces the key of an element.
func (heap *LeftistHeap[E, K]) DecreaseKey(handle base.Handle, key K) {
	node, ok := handle.(*leftistNode[E, K])
	if !ok {
		panic("pq: invalid handle type")
	}
	if cmp.Compare(node.key, key) < 0 {
		panic("pq: DecreaseKey called with key greater than current key")
	}
	heap.Remove(node)
	// we need to add the very same node again because the
	// node's pointer is the handle that should not change!
	node.key = key
	node.up, node.left, node.right = nil, nil, nil
	node.dist = 1
	heap.min = heap.min.merge(node)
	heap.size++
}

// Remove removes an element (by handle) from the heap,
// splitting it in three parts, fixing one of them and
// merging them to one tree again.
func (heap *LeftistHeap[E, K]) Remove(handle base.Handle) {
	node, ok := handle.(*leftistNode[E, K])
	if !ok {
		panic("pq: invalid handle type")
	}
	if node == heap.min {
		heap.ExtractMin()
		return
	}

	lower := node.left.merge(node.right)

	fix := node.up
	if fix.left == node {
		fix.left = nil
	} else {
		fix.right = nil
	}
	for fix != nil {
		if fix.right.distOrZero() > fix.left.distOrZero() {
			fix.left, fix.right = fix.right, fix.left
		}
		if fix.dist > fix.right.distOrZero()+1 {
			fix.dist = fix.right.distOrZero() + 1
		} else {
			break
		}
		fix = fix.up
	}

	heap.min = heap.min.merge(lower)
	heap.size--
}

func (node *LeftistHeap[E, K]) Merge(other *LeftistHeap[E, K]) {
	node.min = node.min.merge(other.min)
	node.size += other.size
	other.min = nil
	other.size = 0
}

func (heap *LeftistHeap[E, K]) MultiInsert(elements []E, keys []K) (handles []base.Handle) {
	if len(elements) == 0 {
		return
	}
	if len(elements) != len(keys) {
		panic("pq: elements and keys must have the same length")
	}
	var nodes []*leftistNode[E, K]
	for i, element := range elements { // O(k)
		node := &leftistNode[E, K]{element: element, key: keys[i], dist: 1}
		nodes = append(nodes, node)
		handles = append(handles, node)
	}
	for len(nodes) > 1 { // O(k)
		newNode := nodes[0].merge(nodes[1])
		nodes = append(nodes, newNode)
		nodes = nodes[2:]
	}
	heap.min = heap.min.merge(nodes[0]) // O(log n)
	heap.size += len(elements)
	return
}

func BuildLeftistHeap[E any, K cmp.Ordered](elements []E, keys []K) (*LeftistHeap[E, K], []base.Handle) {
	res := NewLeftistHeap[E, K]()
	handles := res.MultiInsert(elements, keys)
	return res, handles
}
