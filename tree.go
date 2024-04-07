package goshell

import (
	"fmt"
	"reflect"
	"strings"
)

type tree struct {
	Node  *node
	Count int
}

type node struct {
	Prompt string
	Flag   string
	Params []string
	Data   func(c *Context)
	Tree   *tree
	Left   *node
	Right  *node
}

func newTree() *tree {
	return &tree{
		Count: 0,
	}
}

func newNode(prompt string) *node {
	return &node{
		Prompt: prompt,
		Tree:   newTree(),
	}
}

func (t *tree) AddNode(prompt string, data func(c *Context), flag string, params []string) {
	prompts := strings.Split(prompt, " ")

	t.addHelp()

	var node *node
	if t.Node != nil {
		var inserted bool
		if len(prompts) == 1 {
			node, inserted = t.Node.Insert(prompts[0], params, nil)
		} else {
			node, inserted = t.Node.Insert(prompts[0], nil, nil)
		}
		if inserted {
			t.Count++
		}
	} else {
		t.Node = newNode(prompts[0])
		t.Count++
		node = t.Node
	}

	if len(prompts) > 1 {
		if node.Flag == "" {
			node.Flag = "<Command>"
		}

		pAux := strings.Join(prompts[1:], " ")
		node.Tree.AddNode(pAux, data, flag, params)
		return
	}

	node.Data = data
	node.Flag = flag
	node.Params = params
}

func (t *tree) addHelp() {
	if t.Count == 0 {
		t.Node = &node{
			Prompt: "help",
			Flag:   "Display the commands",
			Tree:   newTree(),
		}
		t.Node.Data = func(c *Context) {
			t.PrintNode()
		}
		t.Count++
	}
}

func (n *node) Insert(prompt string, params []string, nodeP **node) (*node, bool) {
	if n.Prompt == prompt && reflect.DeepEqual(n.Params, params) {
		return n, false
	}

	if n.Prompt == prompt && len(params) == 0 {
		nodeAux := newNode(prompt)
		nodeAux.Left = n
		nodeAux.Right = n.Right
		n.Right = nil
		*nodeP = nodeAux
		return nodeAux, true
	}

	if n.Prompt >= prompt {
		if n.Left == nil {
			n.Left = newNode(prompt)
			return n.Left, true
		}
		return n.Left.Insert(prompt, params, &n.Left)
	} else {
		if n.Right == nil {
			n.Right = newNode(prompt)
			return n.Right, true
		}
		return n.Right.Insert(prompt, params, &n.Right)
	}
}

func (t *tree) FindNode(prompt string) (*node, []string) {
	if t.Count == 0 {
		return nil, nil
	}

	prompts := strings.Split(prompt, " ")

	node := t.Node.Find(prompts[0], false)
	if node != nil && len(prompts) > 1 {
		pAux := strings.Join(prompts[1:], " ")
		nodeAux, params := node.Tree.FindNode(pAux)

		if nodeAux != nil {
			return nodeAux, params
		}
	}
	node = node.paramsChecks(prompts)
	return node, prompts[1:]
}

func (n *node) Find(prompt string, noThis bool) *node {
	if n.Prompt == prompt && !noThis {
		return n
	}

	if n.Prompt >= prompt {
		if n.Left != nil {
			return n.Left.Find(prompt, false)
		}
	}

	if n.Prompt < prompt {
		if n.Right != nil {
			return n.Right.Find(prompt, false)
		}
	}

	return nil
}

func (n *node) paramsChecks(prompts []string) *node {
	if n == nil {
		return n
	}

	if len(n.Params) == len(prompts[1:]) {
		return n
	}
	node := n.Find(prompts[0], true)
	if node != nil {
		if node.Tree.Count > 0 && len(prompts) > 1 {
			pAux := strings.Join(prompts[1:], " ")
			nodeAux, _ := node.Tree.FindNode(pAux)
			if nodeAux != nil {
				node = nodeAux
			}
		}
		return node.paramsChecks(prompts)
	}
	return node
}

func (t *tree) PrintNode() {
	fmt.Println("Usage: ")
	t.Node.print(0)
}

func (n *node) print(repeat int) {
	if n != nil {
		n.Left.print(repeat)
		space := strings.Repeat(" ", repeat)
		prompt := fmt.Sprintf("%-*s", 10, n.Prompt)
		fmt.Printf("\t%s%s %s\n", space, prompt, n.Flag)
		n.Right.print(repeat)
	}
}
