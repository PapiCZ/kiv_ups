package tree

import (
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp/protocol"
	"time"
)

type Node struct {
	Parent   *Node    `json:"-"`
	Children []*Node  `json:"children"`
	Value    GameNode `json:"value"`
}

func NewNode(parent *Node, value GameNode) Node {
	return Node{
		Parent:   parent,
		Children: make([]*Node, 0),
		Value:    value,
	}
}

func (n *Node) GetRoot() *Node {
	node := n

	for n.Parent != nil {
		node = node.Parent
	}

	return node
}

func (n *Node) AddChildren(children ...*Node) {
	for _, child := range children {
		child.Parent = n
		n.Children = append(n.Children, child)
	}
}

func (n *Node) AddGameNodes(gameNodes ...GameNode) {
	for _, gameNode := range gameNodes {
		newNode := NewNode(n, gameNode)
		newNode.Parent = n
		n.Children = append(n.Children, &newNode)
	}
}

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

type GameNode interface {
	Process(playerMessages []interfaces.PlayerMessage, delta time.Duration) // Called every tick
	ListenMessages() []protocol.Message
	Filter(playerMessages []interfaces.PlayerMessage) []interfaces.PlayerMessage
}
