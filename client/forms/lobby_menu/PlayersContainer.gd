extends HBoxContainer

var LobbyPlayer = preload("res://forms/lobby_menu/LobbyPlayer.tscn")

var players = []

func add_player(player_name):	
	var lobby_player = LobbyPlayer.instance()
	lobby_player.get_node("PlayerName").text = player_name

	add_child(lobby_player)
	players.append({
		"name": player_name,
		"node": lobby_player # We need to know this for player removal
	})

func remove_player(player_name):
	for player in players:
		if player_name == player.name:
			players.erase(player)
			player.node.queue_free()
			break