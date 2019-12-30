extends VBoxContainer

var lobby_name
var InGame = preload("res://InGame.tscn")

func _load():
	Network.send({}, MessageTypes.LIST_LOBBY_PLAYERS, self, "_on_lobby_players_loaded")

	Network.connect_message(MessageTypes.START_GAME_RESPONSE, self, "_start_game")
	Network.connect_message(MessageTypes.LOBBY_PLAYER_CONNECTED, self, "_player_connected")

func _on_lobby_players_loaded(data):
	if data[0].response.status:
		for player_name in data[0].response.data.players:
			$FormContainer/PlayersContainer.add_player(player_name)

func _on_Start_pressed():
	Network.disconnect_message(MessageTypes.START_GAME_RESPONSE)
	Network.send({}, MessageTypes.START_GAME, self, "_start_game")

func _start_game(data):
	Network.disconnect_message(MessageTypes.START_GAME_RESPONSE)
	Network.disconnect_message(MessageTypes.LOBBY_PLAYER_CONNECTED)

	Menu.hide_current()
	
	var ingame = InGame.instance()
	ingame.set_name("InGame")
	Network.connect_message(MessageTypes.UPDATE_STATE, ingame, "_update_state")
	Network.connect_message(MessageTypes.GAME_END, self, "_end_game")
	Network.connect_message(MessageTypes.PLAYER_DISCONNECTED, ingame, "_player_disconnected")
	get_parent().get_parent().get_parent().add_child(ingame)

func _player_connected(data):
	$FormContainer/PlayersContainer.add_player(data[0].response.data.name)

func _on_Back_pressed():
	Network.send({
		"name": lobby_name
	}, MessageTypes.DELETE_LOBBY)
	Menu.back()

func _end_game(data):
	Network.disconnect_message(MessageTypes.UPDATE_STATE)
	Network.disconnect_message(MessageTypes.GAME_END)
	Network.disconnect_message(MessageTypes.PLAYER_DISCONNECTED)

	get_tree().get_root().find_node("InGame", true, false).queue_free()
	Menu.reset_all()
	Menu.get(Menu.MENU_LEVEL.END_GAME, false).score_summary = data[0].response.data.score_summary
	Menu.go(Menu.MENU_LEVEL.END_GAME)