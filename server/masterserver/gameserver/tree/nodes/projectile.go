package nodes

import (
	"kiv_ups_server/masterserver/gameserver/tree"
	"kiv_ups_server/masterserver/interfaces"
	"kiv_ups_server/net/tcp/protocol"
)

type Projectile struct {
	PosX      float64           `json:"pos_x"`
	PosY      float64           `json:"pos_y"`
	VelocityX float64           `json:"velocity_x"`
	VelocityY float64           `json:"velocity_y"`
	Rotation  float64           `json:"rotation"`
	Player    interfaces.Player `json:"-"`
	Node      *tree.Node        `json:"-"`
}

func (p *Projectile) Init(node *tree.Node) {
	p.Node = node
}

func (p *Projectile) Process(playerMessages []interfaces.PlayerMessage, delta float64) {
	p.PosX += p.VelocityX * delta
	p.PosY += p.VelocityY * delta

	if p.PosX < 0 || p.PosX > 1920 || p.PosY < 0 || p.PosY > 1080 {
		p.Node.Destroy()
		return
	}

	// Check collision with asteroid
	for _, node := range p.Node.GetRoot().FindAllChildrenByType("asteroid") {
		asteroid := node.Value.(*Asteroid)

		asteroidCollider := Circle{
			asteroid.PosX,
			asteroid.PosY,
			asteroid.Radius,
		}

		if asteroidCollider.IsPointInside(p.PosX, p.PosY) {
			node.Destroy()
			p.Node.Destroy()
			p.AddPlayerScore(p.Player, asteroid.Value)

			if asteroid.Scale > 0.4 {
				// Create 2 new asteroids with rotated velocity about 0.3 and -0.3rad
				// of origin asteroid
				v := Vector{asteroid.VelocityX, asteroid.VelocityY}

				newAsteroid1 := Asteroid{
					PosX:      asteroid.PosX,
					PosY:      asteroid.PosY,
					VelocityX: v.Rotated(0.3).X,
					VelocityY: v.Rotated(0.3).Y,
					Rotation:  asteroid.Rotation,
					Scale:     asteroid.Scale / 2.0,
					Value:     asteroid.Value / 2.0,
					Node:      nil,
				}

				newAsteroid2 := Asteroid{
					PosX:      asteroid.PosX,
					PosY:      asteroid.PosY,
					VelocityX: v.Rotated(-0.3).X,
					VelocityY: v.Rotated(-0.3).Y,
					Rotation:  asteroid.Rotation,
					Scale:     asteroid.Scale / 2.0,
					Value:     asteroid.Value / 2.0,
					Node:      nil,
				}

				node1 := tree.NewNode(node.Parent, &newAsteroid1)
				node1.Init()
				node1.Value.Init(&node1)
				node2 := tree.NewNode(node.Parent, &newAsteroid2)
				node2.Init()
				node2.Value.Init(&node2)
				node.Parent.Children = append(node.Parent.Children, &node1)
				node.Parent.Children = append(node.Parent.Children, &node2)
			}

			break
		}
	}

	// Check collision with enemy spaceship
	for _, node := range p.Node.GetRoot().FindAllChildrenByType("spaceship") {
		spaceship := node.Value.(*Spaceship)

		if spaceship.Player == p.Player {
			// Player cant kill itself
			continue
		}

		spaceshipCollider := Circle{
			spaceship.PosX,
			spaceship.PosY,
			spaceship.Radius,
		}

		if spaceshipCollider.IsPointInside(p.PosX, p.PosY) {
			spaceship.Die()
			p.AddPlayerScore(p.Player, 1000)

			break
		}
	}

}

// AddPlayerScore adds given amount of score to given player
func (p *Projectile) AddPlayerScore(player interfaces.Player, score int) {
	for _, node := range p.Node.GetRoot().FindAllChildrenByType("score") {
		scoreGameNode := node.Value.(*Score)

		if scoreGameNode.Player == player {
			scoreGameNode.Score += score
			return
		}
	}
}

func (p *Projectile) ListenMessages() []protocol.Message {
	return []protocol.Message{}
}

func (p *Projectile) Filter(playerMessages []interfaces.PlayerMessage) []interfaces.PlayerMessage {
	return playerMessages
}
