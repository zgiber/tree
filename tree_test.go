package tree

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func newTestTree() *Tree {

	tree := NewTree(NewNode("root"))
	for i := 0; i < 10; i++ {
		tree.root.AddNewChild(
			strings.Join(
				[]string{"user", strconv.Itoa(i)}, ""))
	}

	for _, cn := range tree.root.children {
		for i := 0; i < 10; i++ {
			cn.AddNewChild(
				strings.Join(
					[]string{"hub", strconv.Itoa(i)}, ""))
		}
	}

	return tree
}

func TestTree(t *testing.T) {
	tree := newTestTree()

	// creating a node with SetPath
	n := NewNode("hub80")
	tree.root.SetChildByPath(n, "user19")
	err := tree.SetNodeStatus(map[string]interface{}{"online": true}, false, false, "user19", "hub80")
	if err != nil {
		t.Fatal(err)
	}

	// retrieving a node's status with GetStatus
	status, err := tree.NodeStatus("user19", "hub80")
	if err != nil {
		t.Fatal(err)
	}

	// checking if status was set
	if _, ok := status["online"]; !ok {
		t.Fatal("Status is not set.")
	}

	// check if the other node additions were correctly performed
	if len(tree.root.children) != 11 {
		t.Fatal("Root is expected to have 11 children.")
	}

	// retrieve an arbitrary node
	_, err = tree.root.ChildByPath("user2", "hub1")
	if err != nil {
		t.Fatal(err)
	}

	// retrieving the node we've created manually
	n, err = tree.root.ChildByPath("user19", "hub80")
	if err != nil {
		t.Fatal(err)
	}

}

func TestSettingsRecursive(t *testing.T) {
	tree := newTestTree()

	err := tree.SetNodeStatus(map[string]interface{}{"online": false}, true, false, "user2")
	if err != nil {
		t.Fatal(err)
	}

	status, err := tree.NodeStatus("user2", "hub3")
	if err != nil {
		t.Fatal(err)
	}

	if value, ok := status["online"]; !ok {
		t.Fatal("Recursive status set failed, not found.")
	} else {
		if online, ok := value.(bool); !ok {
			t.Fatal("Recursive status set failed, wrong type.")
		} else {
			if online != false {
				t.Fatal("Recursive status set failed, wrong value.")
			}
		}
	}
}

func TestSettingsBubbleUp(t *testing.T) {
	tree := newTestTree()

	err := tree.SetNodeStatus(map[string]interface{}{"online": true}, false, true, "user3", "hub3")
	if err != nil {
		t.Fatal(err)
	}

	status, err := tree.NodeStatus("user3")
	if err != nil {
		t.Fatal(err)
	}

	if value, ok := status["online"]; !ok {
		t.Fatal("Recursive status set failed, not found.")
	} else {
		if online, ok := value.(bool); !ok {
			t.Fatal("Recursive status set failed, wrong type.")
		} else {
			if online != true {
				t.Fatal("Recursive status set failed, wrong value.")
			}
		}
	}
}

func TestValue(t *testing.T) {
	tree := newTestTree()
	n, err := tree.Node("user2", "hub3")
	if err != nil {
		t.Fatal(err)
	}

	//
	// Float value
	//

	err = n.SetValue(34.56)
	if err != nil {
		t.Fatal(err)
	}

	if vf, err := n.value.Float(); err != nil {
		t.Fatal(err)
	} else if vf != 34.56 {
		t.Fatal("value mismatch")
	}

	b, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	//
	// Bool value
	//

	err = n.SetValue(true)
	if err != nil {
		t.Fatal(err)
	}

	if vb, err := n.value.Bool(); err != nil {
		t.Fatal(err)
	} else if vb != true {
		t.Fatal("value mismatch")
	}

	b, err = json.MarshalIndent(n, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	//
	// Map value
	//

	err = n.SetValue(map[string]interface{}{"x": 12.34, "y": true})
	if err != nil {
		t.Fatal(err)
	}

	if vm, err := n.value.Map(); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(vm)
	}

	b, err = json.MarshalIndent(n, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
	fmt.Println(tree)
}
