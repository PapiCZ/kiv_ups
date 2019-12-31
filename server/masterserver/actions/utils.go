package actions

import (
	"kiv_ups_server/masterserver/interfaces"
	"math/rand"
)

// ConvertShadowPlayerToPlayer gives to player valid UID and name
func ConvertShadowPlayerToPlayer(player interfaces.Player, name string) interfaces.Player {
	player.SetUID(interfaces.PlayerUID(rand.Int()))
	player.SetName(name)

	return player
}
