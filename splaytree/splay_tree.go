package splaytree

import (
	"cmp"
	"fmt"
)

type Node[K cmp.Ordered, V any] struct {
	Key    K
	Value  V
	left   *Node[K, V]
	right  *Node[K, V]
	parent *Node[K, V]
}

type SplayTree[K cmp.Ordered, V any] struct {
	root        *Node[K, V]
	size        int
	OnSplayStep func(xKey K, stepDescription string)
}

func New[K cmp.Ordered, V any]() *SplayTree[K, V] {
	return &SplayTree[K, V]{}
}

func (t *SplayTree[K, V]) Size() int { return t.size }

func (t *SplayTree[K, V]) Root() (K, bool) {
	if t.root == nil {
		var zero K
		return zero, false
	}
	return t.root.Key, true
}

func (t *SplayTree[K, V]) rotateRight(x *Node[K, V]) {
	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		t.root = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}
	y.right = x
	x.parent = y
}

func (t *SplayTree[K, V]) rotateLeft(x *Node[K, V]) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (t *SplayTree[K, V]) splay(x *Node[K, V]) {
	for x.parent != nil {
		p := x.parent
		g := p.parent
		var stepDesc string

		if g == nil {
			if x == p.left {
				t.rotateRight(p)
				stepDesc = fmt.Sprintf("Zig Derecho (rotación derecha sobre el padre %v)", p.Key)
			} else {
				t.rotateLeft(p)
				stepDesc = fmt.Sprintf("Zig Izquierdo (rotación izquierda sobre el padre %v)", p.Key)
			}
		} else if x == p.left && p == g.left {
			t.rotateRight(g)
			t.rotateRight(p)
			stepDesc = fmt.Sprintf("Zig-Zig Derecho (rotación derecha sobre el abuelo %v y luego sobre el padre %v)", g.Key, p.Key)
		} else if x == p.right && p == g.right {
			t.rotateLeft(g)
			t.rotateLeft(p)
			stepDesc = fmt.Sprintf("Zig-Zig Izquierdo (rotación izquierda sobre el abuelo %v y luego sobre el padre %v)", g.Key, p.Key)
		} else if x == p.right && p == g.left {
			t.rotateLeft(p)
			t.rotateRight(g)
			stepDesc = fmt.Sprintf("Zig-Zag Izquierdo-Derecho (rotación izquierda sobre el padre %v y luego derecha sobre el abuelo %v)", p.Key, g.Key)
		} else {
			t.rotateRight(p)
			t.rotateLeft(g)
			stepDesc = fmt.Sprintf("Zig-Zag Derecho-Izquierdo (rotación derecha sobre el padre %v y luego izquierda sobre el abuelo %v)", p.Key, g.Key)
		}

		if t.OnSplayStep != nil {
			t.OnSplayStep(x.Key, stepDesc)
		}
	}
}

func (t *SplayTree[K, V]) Insert(key K, value V) {
	if t.root == nil {
		t.root = &Node[K, V]{Key: key, Value: value}
		t.size++
		return
	}

	cur := t.root
	var parent *Node[K, V]
	for cur != nil {
		parent = cur
		switch {
		case key < cur.Key:
			cur = cur.left
		case key > cur.Key:
			cur = cur.right
		default:
			cur.Value = value
			t.splay(cur)
			return
		}
	}

	n := &Node[K, V]{Key: key, Value: value, parent: parent}
	if key < parent.Key {
		parent.left = n
	} else {
		parent.right = n
	}
	t.size++
	t.splay(n)
}

func (t *SplayTree[K, V]) Search(key K) (V, bool) {
	cur := t.root
	var last *Node[K, V]
	for cur != nil {
		last = cur
		switch {
		case key < cur.Key:
			cur = cur.left
		case key > cur.Key:
			cur = cur.right
		default:
			t.splay(cur)
			return cur.Value, true
		}
	}
	if last != nil {
		t.splay(last)
	}
	var zero V
	return zero, false
}

func (t *SplayTree[K, V]) Delete(key K) bool {
	_, found := t.Search(key)
	if !found {
		return false
	}

	left := t.root.left
	right := t.root.right
	t.size--

	if left == nil {
		t.root = right
		if right != nil {
			right.parent = nil
		}
		return true
	}
	if right == nil {
		t.root = left
		left.parent = nil
		return true
	}

	left.parent = nil
	t.root = left
	cur := left
	for cur.right != nil {
		cur = cur.right
	}
	t.splay(cur)

	t.root.right = right
	right.parent = t.root
	return true
}

func (t *SplayTree[K, V]) InOrder() []K {
	result := make([]K, 0, t.size)
	var traverse func(n *Node[K, V])
	traverse = func(n *Node[K, V]) {
		if n == nil {
			return
		}
		traverse(n.left)
		result = append(result, n.Key)
		traverse(n.right)
	}
	traverse(t.root)
	return result
}

func (t *SplayTree[K, V]) Height() int {
	var h func(n *Node[K, V]) int
	h = func(n *Node[K, V]) int {
		if n == nil {
			return 0
		}
		lh := h(n.left)
		rh := h(n.right)
		if lh > rh {
			return lh + 1
		}
		return rh + 1
	}
	return h(t.root)
}

func (n *Node[K, V]) Left() *Node[K, V] {
	if n == nil {
		return nil
	}
	return n.left
}

func (n *Node[K, V]) Right() *Node[K, V] {
	if n == nil {
		return nil
	}
	return n.right
}

func (t *SplayTree[K, V]) RootNode() *Node[K, V] {
	return t.root
}


