package tree

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

type Node struct {
	sync.RWMutex
	id       string
	tree     *Tree
	value    *Value
	children map[string]*Node
	parent   *Node
	status   map[string]interface{}
}

func NewNode(id string) *Node {
	return &Node{
		id:       id,
		children: map[string]*Node{},
		status:   map[string]interface{}{},
	}
}

func (n *Node) SetValue(v interface{}) error {
	value, err := NewValue(v)
	if err != nil {
		return err
	}

	n.value = value
	return nil
}

func (n *Node) AddNewChild(id string) {
	nn := NewNode(id)
	nn.tree = n.tree
	nn.parent = n
	n.Lock()
	n.children[id] = nn
	n.Unlock()
}

func (n *Node) AddChild(nn *Node) {
	n.Lock()
	nn.tree = n.tree
	nn.parent = n
	n.children[nn.id] = nn
	n.Unlock()
}

func (n *Node) Child(id string) *Node {
	n.RLock()
	defer n.RUnlock()
	return n.children[id]
}

func (n *Node) Value() *Value {
	return n.value
}

func (n *Node) Map() map[string]interface{} {
	m := map[string]interface{}{}
	m["children"] = map[string]interface{}{}
	m["status"] = n.status

	n.RLock()
	defer n.RUnlock()
	for id, cn := range n.children {
		m["children"].(map[string]interface{})[id] = cn.Map()
	}

	return m
}

func (n *Node) ChildByPath(path ...string) (*Node, error) {
	if len(path) == 0 {
		return n, nil
	}

	nn, ok := n.children[path[0]]
	if ok {
		return nn.ChildByPath(path[1:]...)
	}

	return nil, fmt.Errorf("node not found")
}

func (n *Node) SetChildByPath(nn *Node, path ...string) {
	if len(path) == 0 {
		n.AddChild(nn)
		return
	}

	n.RLock()
	cn, ok := n.children[path[0]]
	n.RUnlock()
	if !ok {
		cn = NewNode(path[0])
	}

	cn.SetChildByPath(nn, path[1:]...)
	n.AddChild(cn)
}

func (n *Node) DeletePath(path ...string) {
	if len(path) == 1 {
		n.Lock()
		for _, cn := range n.children {
			cn.tree = nil   // remove reference to the tree
			cn.parent = nil // and parents
		}
		delete(n.children, path[0])
		n.Unlock()
	}
}

func (n *Node) Version() (uint64, error) {
	if n.tree == nil {
		return 0, fmt.Errorf("node is not part of a tree")
	} else {
		return n.tree.Version(), nil
	}
}

//func (n *Node) cut() {
//	if len(n.children) > 0 {
//		for _, cn := range n.children {
//			cn.cut()
//		}
//	} else {
//		n.tree = nil
//	}
//}

func (n *Node) SetStatus(values map[string]interface{}, recursive bool) {
	if values == nil {
		return
	}
	n.Lock()
	for k, v := range values {
		n.status[k] = v
	}
	n.Unlock()

	if recursive {
		n.RLock()
		for _, cn := range n.children {
			cn.SetStatus(values, true)
		}
		n.RUnlock()
	}
}

func (n *Node) Status() map[string]interface{} {
	// copy values
	m := map[string]interface{}{}
	n.RLock()
	for k, v := range n.status {
		m[k] = v
	}
	n.RUnlock()
	return m
}

func (n *Node) MarshalJSON() ([]byte, error) {
	type proxy struct {
		ID       string                 `json:"id"`
		Value    *Value                 `json:"value,omitempty"`
		Status   map[string]interface{} `json:"status"`
		Children map[string]*Node       `json:"children"`
	}

	p := &proxy{
		ID:       n.id,
		Value:    n.value,
		Status:   n.status,
		Children: n.children,
	}

	return json.Marshal(p)
}

// UnmarshalJSON constructs a Node from it's JSON encoded value.
// A Node created this way is not necessarily ready to be used
// on a Tree.
func (n *Node) UnmarshalJSON(b []byte) error {
	type proxy struct {
		ID       string                 `json:"id"`
		Value    *Value                 `json:"value,omitempty"`
		Status   map[string]interface{} `json:"status"`
		Children map[string]*Node       `json:"children"`
	}

	p := &proxy{}

	dec := json.NewDecoder(bytes.NewBuffer(b))
	dec.UseNumber()
	err := dec.Decode(p)
	if err != nil {
		return err
	}

	n.id = p.ID
	n.value = p.Value
	n.status = p.Status
	n.children = p.Children

	return nil
}

func (n *Node) String() string {
	out, _ := json.MarshalIndent(n, "", "   ")
	return string(out)
}
