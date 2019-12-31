package actions

import (
	interfaces2 "kiv_ups_server/internal/masterserver/interfaces"
	"math/rand"
)

// ConvertShadowPlayerToPlayer gives to player valid UID and name
func ConvertShadowPlayerToPlayer(player interfaces2.Player, name string) interfaces2.Player {
	player.SetUID(interfaces2.PlayerUID(rand.Int()))
	player.SetName(name)

	return player
}
