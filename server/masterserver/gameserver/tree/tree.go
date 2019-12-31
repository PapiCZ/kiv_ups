package tree

import (
	"kiv_ups_server/masterserver/interfaces"
	"kiv_ups_server/net/tcp/protocol"
	"math/rand"
	"reflect"
	"strings"
)

const NodeRandomIdLength = 10
const NodeRandomIdValidChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Node is basic element of game tree. It is used to handle all game objects in
// a tree structure.
type Node struct {
	Parent   *Node    `json:"-"`
	Children []*Node  `json:"children"`
	Value    GameNode `json:"value"`
	Type     string   `json:"type"`
	Id       string   `json:"id"`
}

// NewNode creates and initializes new node with given value.
func NewNode(parent *Node, value GameNode) Node {
	return Node{
		Parent:   parent,
		Children: make([]*Node, 0),
		Value:    value,
	}
}

// Init initializes node's type and unique identifier
func (n *Node) Init() {
	if n.Value != nil {
		n.Type = strings.ToLower(reflect.TypeOf(n.Value).Elem().Name())
	}

	n.Id = n.Type + "_" + RandomString(
		[]rune(NodeRandomIdValidChars),
		NodeRandomIdLength,
	)
}

// GetRoot returns root of tree (node that doesn't have parent)
func (n *Node) GetRoot() *Node {
	node := n

	for node.Parent != nil {
		node = node.Parent
	}

	return node
}

// AddChildren adds given children to node
func (n *Node) AddChildren(children ...*Node) {
	for _, child := range children {
		child.Parent = n
		n.Children = append(n.Children, child)
	}
}

// AddGameNodes adds game nodes to node
func (n *Node) AddGameNodes(gameNodes ...GameNode) {
	for _, gameNode := range gameNodes {
		newNode := NewNode(n, gameNode)
		newNode.Parent = n
		n.Children = append(n.Children, &newNode)
	}
}

// GetAllChildren returns all children of tree (recursively)
func (n *Node) GetAllChildren() []*Node {
	nodes := make([]*Node, 0)
	nodes = append(nodes, n.Children...)

	if n.Children != nil && len(n.Children) > 0 {
		for _, child := range n.Children {
			nodes = append(nodes, child.GetAllChildren()...)
		}
	}

	return nodes
}

// FindAllChildrenByType finds all child nodes in the tree
// by given type (recursively)
func (n *Node) FindAllChildrenByType(type_ string) []*Node {
	out := make([]*Node, 0)
	for _, node := range n.GetAllChildren() {
		if node.Type == type_ {
			out = append(out, node)
		}
	}
	return out
}

// Destroy removes itself from parent node
func (n *Node) Destroy() {
	// find node in parent node
	if n.Parent != nil {
		for i, node := range n.Parent.Children {
			if node == n {
				n.Parent.Children = append(n.Parent.Children[:i], n.Parent.Children[i+1:]...)
				break
			}
		}
	}
}

// RandomString build string from randomly chosen characters of given length
func RandomString(charset []rune, length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GameNode is interface for all in-game objects. Function Process is called
// every game tick. Function ListenMessages is used to tell game server for
// what messages object want to listen and function Filter can be used to
// add custom filtering given message instances.
type GameNode interface {
	Init(node *Node)
	Process(playerMessages []interfaces.PlayerMessage, delta float64) // Called every tick
	ListenMessages() []protocol.Message
	Filter(playerMessages []interfaces.PlayerMessage) []interfaces.PlayerMessage
}
