package nodes

import (
	log "github.com/sirupsen/logrus"
	"kiv_ups_server/masterserver/gameserver/settings"
	"kiv_ups_server/masterserver/gameserver/tree"
	"kiv_ups_server/masterserver/interfaces"
	"kiv_ups_server/net/tcp/protocol"
	"math"
	"math/rand"
)

const ShootTimeout = 0.15
const Friction = 0.7
const MaxSpeed = 500
const RotationSpeed = 200
const AntiCheatTolerance = 0.1
const AntiCheatPositionTolerance = 20

type Spaceship struct {
	PosX                   float64           `json:"pos_x"`
	PosY                   float64           `json:"pos_y"`
	VelocityX              float64           `json:"velocity_x"`
	VelocityY              float64           `json:"velocity_y"`
	Rotation               float64           `json:"rotation"`
	PlayerName             string            `json:"player_name"`
	Immune                 bool              `json:"immune"`
	ReloadPosition         bool              `json:"reload_position"`
	TeleportAllowedCounter int               `json:"-"` // 0 = teleport is not allowed
	ShootTimeout           float64           `json:"-"`
	Player                 interfaces.Player `json:"-"`
	Score                  *Score            `json:"-"`
	Speed                  float64           `json:"-"`
	ImmuneTime             float64           `json:"-"`
	Radius                 float64           `json:"-"`
	Node                   *tree.Node        `json:"-"`
}

func (s *Spaceship) Init(node *tree.Node) {
	s.Node = node
	s.Radius = 50

	if s.Player != nil {
		s.PlayerName = s.Player.GetName()
	}
}

func (s *Spaceship) Process(playerMessages []interfaces.PlayerMessage, delta float64) {
	if s.ReloadPosition {
		s.ReloadPosition = false
	}
	if s.ShootTimeout >= 0 {
		// Decrement shoot timeout
		s.ShootTimeout -= delta
	}
	if s.TeleportAllowedCounter > 0 {
		s.TeleportAllowedCounter -= 1
	}

	for _, playerMessage := range playerMessages {
		message := playerMessage.GetMessage().Message
		if v, ok := message.(*protocol.PlayerMoveMessage); ok {
			s.AntiCheat(v, delta)

			s.PosX = v.PosX
			s.PosY = v.PosY
			s.VelocityX = v.VelocityX
			s.VelocityY = v.VelocityY
			s.Rotation = v.Rotation
			s.PlayerName = s.Player.GetName()
		} else if _, ok := message.(*protocol.ShootProjectileMessage); ok {
			// Check if player can shoot
			if s.ShootTimeout <= 0 && !s.Immune {
				s.ShootTimeout = ShootTimeout

				velocity := Vector{0, -1}
				velocity = velocity.Rotated(s.Rotation).Normalized()

				node := tree.NewNode(s.Node.Parent, &Projectile{
					PosX:      s.PosX,
					PosY:      s.PosY,
					VelocityX: velocity.X * 600,
					VelocityY: velocity.Y * 600,
					Rotation:  s.Rotation,
					Player:    s.Player,
					Node:      nil,
				})
				node.Init()
				node.Value.Init(&node)

				s.Node.Parent.Children = append(s.Node.Parent.Children, &node)
			}
		}
	}

	// check collision with asteroid
	if !s.Immune {
		for _, node := range s.Node.GetRoot().FindAllChildrenByType("asteroid") {
			asteroid := node.Value.(*Asteroid)
			asteroidCollider := Circle{
				asteroid.PosX,
				asteroid.PosY,
				asteroid.Radius,
			}

			if asteroidCollider.IsPointInside(s.PosX, s.PosY) {
				s.Die()
			}
		}
	}

	// remove immunity
	if s.Immune {
		s.ImmuneTime += delta
		if s.ImmuneTime > 3 {
			s.Immune = false
		}
	}
}

func (s *Spaceship) ListenMessages() []protocol.Message {
	return []protocol.Message{protocol.PlayerMoveMessage{}, protocol.ShootProjectileMessage{}}
}

func (s *Spaceship) Filter(playerMessages []interfaces.PlayerMessage) []interfaces.PlayerMessage {
	filteredPlayerMessages := make([]interfaces.PlayerMessage, 0)

	for _, playerMessage := range playerMessages {
		if _, ok := playerMessage.GetMessage().Message.(*protocol.PlayerMoveMessage); ok &&
			playerMessage.GetPlayer().GetName() != s.Player.GetName() {
			continue
		}

		if _, ok := playerMessage.GetMessage().Message.(*protocol.ShootProjectileMessage); ok &&
			playerMessage.GetPlayer().GetName() != s.Player.GetName() {
			continue
		}

		filteredPlayerMessages = append(filteredPlayerMessages, playerMessage)
	}

	return filteredPlayerMessages
}

func (s *Spaceship) RandomizePosition() {
	s.PosX = float64(rand.Intn(settings.Width))
	s.PosY = float64(rand.Intn(settings.Height))
	s.VelocityX = 0
	s.VelocityY = 0
	s.ReloadPosition = true
	s.TeleportAllowedCounter = 2
}

func (s *Spaceship) Immunity() {
	s.Immune = true
	s.ImmuneTime = 0
}

func (s *Spaceship) Die() {
	s.Immunity()
	s.RandomizePosition()
}

func (s *Spaceship) AntiCheat(m *protocol.PlayerMoveMessage, delta float64) {
	// Check if position is in valid range
	if m.PosX < 0 || m.PosX > settings.Width || m.PosY < 0 || m.PosY > settings.Height {
		s.Player.IncrementCheatCounter()
	}

	// Check rotation
	if !IsNumberInTolerance(
		s.Rotation-(DegToRad(RotationSpeed)*delta), m.Rotation, AntiCheatTolerance,
	) && !IsNumberInTolerance(
		s.Rotation+(DegToRad(RotationSpeed)*delta), m.Rotation, AntiCheatTolerance,
	) && !IsNumberInTolerance(s.Rotation, m.Rotation, AntiCheatTolerance) {
		log.Errorln("Invalid rotation")
		s.Player.IncrementCheatCounter()
	}

	// m.Rotation contains valid value, it's safe to use
	// Check velocity angle
	v := Vector{0, -1}
	if m.VelocityX != 0 && m.VelocityY != 0 && !IsNumberInTolerance(
		Vector{m.VelocityX, m.VelocityY}.Rotated(-m.Rotation).Angle(v), 0, AntiCheatTolerance,
	) {
		log.Errorln("Invalid velocity angle")
		s.Player.IncrementCheatCounter()
	}

	// Check speed
	if MaxSpeed < int(math.Abs(
		Vector{s.VelocityX, s.VelocityY}.Length()-Vector{m.VelocityX, m.VelocityY}.Length(),
	)) && !(IsNumberInTolerance(m.PosX, 0, AntiCheatPositionTolerance) ||
		IsNumberInTolerance(m.PosX, settings.Width, AntiCheatPositionTolerance) ||
		IsNumberInTolerance(m.PosY, 0, AntiCheatPositionTolerance) ||
		IsNumberInTolerance(m.PosY, settings.Height, AntiCheatPositionTolerance)) &&
		s.TeleportAllowedCounter == 0 {
		log.Errorln("Invalid speed")
		s.Player.IncrementCheatCounter()
	}

	// m.VelocityX and m.VelocityY contains valid values, it's safe to use
	validArea := Circle{s.PosX, s.PosY, Vector{m.VelocityX, m.VelocityY}.Length() + AntiCheatPositionTolerance}
	if !validArea.IsPointInside(m.PosX, m.PosY) && s.TeleportAllowedCounter == 0 {
		log.Errorln("Invalid position")
		s.Player.IncrementCheatCounter()
	}

	s.Player.ResetCheatCounter()
}

func IsNumberInTolerance(shouldBe float64, number float64, tolerance float64) bool {
	min := shouldBe - tolerance
	max := shouldBe + tolerance

	return min < number && number < max
}
