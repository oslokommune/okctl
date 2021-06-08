package dependencytree

import (
	"bytes"
	"io"
)

// Equals knows how to compare two ResourceNodes and determine equality
func (node *Node) Equals(targetNode *Node) bool {
	if targetNode == nil {
		return false
	}

	// P.S.: For now, this is good enough due to all types of resources only existing once
	return node.Type == targetNode.Type
}

// GetNode returns an identical node as targetNode from the receiver's tree
func (node *Node) GetNode(targetNode *Node) *Node {
	if node.Equals(targetNode) {
		return node
	}

	for _, child := range node.Children {
		result := child.GetNode(targetNode)

		if result != nil {
			return result
		}
	}

	return nil
}

func (node *Node) treeViewGenerator(writer io.Writer, tabs int) {
	result := ""
	for index := 0; index < tabs; index++ {
		result += "\t"
	}

	result += "- " + node.Type.String()

	if node.State == NodeStatePresent {
		result += " (add)" + "\n"
		_, _ = writer.Write([]byte(result))
	}

	if node.State == NodeStateAbsent {
		result += " (remove)" + "\n"
		_, _ = writer.Write([]byte(result))
	}

	if node.State == NodeStateNoop {
		result += " (noop)" + "\n"
		_, _ = writer.Write([]byte(result))
	}

	for _, child := range node.Children {
		child.treeViewGenerator(writer, tabs+1)
	}
}

func (node *Node) String() string {
	var buf bytes.Buffer

	node.treeViewGenerator(&buf, 0)

	return buf.String()
}

// AppendChild adds a child node to the node
func (node *Node) AppendChild(children ...*Node) {
	node.Children = append(node.Children, children...)
}

// NewNode creates a node of the certain type
func NewNode(nodeType NodeType) (child *Node) {
	child = &Node{
		Type:     nodeType,
		Children: make([]*Node, 0),
	}

	child.State = NodeStatePresent

	return child
}

// ApplyFunction will use the supplied ApplyFn on all the nodes in the receiver tree
func (node *Node) ApplyFunction(fn ApplyFn) {
	for _, child := range node.Children {
		child.ApplyFunction(fn)
	}

	fn(node)
}

// ApplyFunctionWithTarget will use the supplied ApplyFnWithTarget on all the nodes in the receiver tree, with an equal
// node from the target tree
func (node *Node) ApplyFunctionWithTarget(fn ApplyFnWithTarget, targetTree *Node) {
	for _, child := range node.Children {
		child.ApplyFunctionWithTarget(fn, targetTree)
	}

	targetNode := targetTree.GetNode(node)
	fn(node, targetNode)
}
