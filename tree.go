package tree

import (
	"encoding/json"
	"sync"
	"sync/atomic"
)

type version struct {
	v *uint64
}

func (v *version) add(delta uint64) {
	atomic.AddUint64(v.v, delta)
}

func (v *version) set(nv uint64) {
	atomic.StoreUint64(v.v, nv)
}

type Tree struct {
	sync.RWMutex
	root    *Node
	version *version
}

func NewTree(rootNode *Node) *Tree {
	var v uint64
	t := &Tree{
		root:    rootNode,
		version: &version{&v},
	}

	rootNode.tree = t
	return t
}

func (tree *Tree) Version() uint64 {
	return atomic.LoadUint64(tree.version.v)
}

func (tree *Tree) SetNodeStatus(values map[string]interface{}, recursive bool, bubbleUp bool, path ...string) error {

	n, err := tree.root.ChildByPath(path...)
	if err != nil {
		// node not found
		// find the existing parent node on the path
		// and create the missing nodes
		n = tree.root
		for i := 0; i < len(path); i++ {
			if nn, ok := n.children[path[i]]; ok {
				n = nn
			} else {
				//log.Printf("creating new node id=%s path=%s\n", id, path[:i+1])

				nn := NewNode(path[i])
				n.AddChild(nn)
				n = nn
			}
		}
	}

	n.SetStatus(values, recursive)

	if bubbleUp {
		for {
			if n.parent != nil {
				n.parent.SetStatus(values, false)
				n = n.parent
			} else {
				break
			}
		}
	}

	tree.version.add(1)
	return nil
}

func (tree *Tree) NodeStatus(path ...string) (map[string]interface{}, error) {
	tree.RLock()
	defer tree.RUnlock()

	n, err := tree.root.ChildByPath(path...)
	if err != nil {
		return nil, err
	}

	n.RLock()
	defer n.RUnlock()
	tree.version.add(1)
	return n.status, nil
}

func (tree *Tree) NewNode(id string, path ...string) *Node {
	nn := NewNode(id)
	tree.root.SetChildByPath(nn, path...)
	tree.version.add(1)
	return nn
}

func (tree *Tree) Node(path ...string) (*Node, error) {
	n, err := tree.root.ChildByPath(path...)
	if err != nil {
		return nil, err
	}
	tree.version.add(1)
	return n, nil
}

func (tree *Tree) DeleteNode(path ...string) error {
	tree.root.DeletePath(path...)
	tree.version.add(1)
	// TODO: if delete did not do anything, don't change version
	return nil
}

func (tree *Tree) String() string {
	out, _ := json.MarshalIndent(tree.root, "", "   ")
	return string(out)
}
