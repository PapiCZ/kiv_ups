extends VBoxContainer

var lobby_name
var InGame = preload("res://InGame.tscn")

func _load():
	Network.send({}, MessageTypes.LIST_LOBBY_PLAYERS, self, "_on_lobby_players_loaded")

	Network.connect_message(MessageTypes.START_GAME_RESPONSE, self, "_start_game")
	Network.connect_message(MessageTypes.LOBBY_PLAYER_CONNECTED, self, "_player_connected")
	Network.connect_message(MessageTypes.LOBBY_PLAYER_DISCONNECTED, self, "_player_disconnected")
	Network.connect("disconnected", self, "_network_disconnected", [], Network.CONNECT_ONESHOT)

func _network_disconnected():
	Network.disconnect_message(MessageTypes.START_GAME_RESPONSE)
	Network.disconnect_message(MessageTypes.LOBBY_PLAYER_CONNECTED)
	Network.disconnect_message(MessageTypes.LOBBY_PLAYER_DISCONNECTED)

	if get_tree().get_root().find_node("InGame", true, false):
		leave()
	else:
		Menu.go(Menu.MENU_LEVEL.MAIN)

func _on_lobby_players_loaded(data):
	# Show connected players
	if data[0].response.status:
		for player_name in data[0].response.data.players:
			$FormContainer/PlayersContainer.add_player(player_name)

func _on_Start_pressed():
	Network.disconnect_message(MessageTypes.START_GAME_RESPONSE)
	Network.send({}, MessageTypes.START_GAME, self, "_start_game")

func _start_game(data):
	# Start the game and disconnect from lobby-related message types
	Network.disconnect_message(MessageTypes.START_GAME_RESPONSE)
	Network.disconnect_message(MessageTypes.LOBBY_PLAYER_CONNECTED)
	Network.disconnect_message(MessageTypes.LOBBY_PLAYER_DISCONNECTED)

	Menu.hide_current()
	
	var ingame = InGame.instance()
	ingame.set_name("InGame")
	# Create listeners for in-game-related messages
	Network.connect_message(MessageTypes.UPDATE_STATE, ingame, "_update_state")
	Network.connect_message(MessageTypes.GAME_END, self, "_end_game")
	Network.connect_message(MessageTypes.PLAYER_DISCONNECTED, ingame, "_player_disconnected")
	Network.connect_message(MessageTypes.PLAYER_CONNECTED, ingame, "_player_connected")
	get_parent().get_parent().get_parent().add_child(ingame)

func _input(ev):
	if Input.is_action_pressed("ui_leave_ingame") and get_tree().get_root().find_node("InGame", true, false):
		leave()

func leave():
		# End the game and show winner
		Network.disconnect_message(MessageTypes.UPDATE_STATE)
		Network.disconnect_message(MessageTypes.GAME_END)
		Network.disconnect_message(MessageTypes.PLAYER_DISCONNECTED)
		Network.disconnect_message(MessageTypes.PLAYER_CONNECTED)

		get_tree().get_root().find_node("InGame", true, false).queue_free()
		Menu.reset_all()
		Menu.go(Menu.MENU_LEVEL.MAIN)

func _player_connected(data):
	# Add palyer to the lobby
	$FormContainer/PlayersContainer.add_player(data[0].response.data.name)

func _player_disconnected(data):
	# Remove palyer from the lobby
	$FormContainer/PlayersContainer.remove_player(data[0].response.data.name)

func _on_Back_pressed():
	Network.send({}, MessageTypes.LEAVE_LOBBY)
	Menu.back()

func _end_game(data):
	# End the game and show winner
	Network.disconnect_message(MessageTypes.UPDATE_STATE)
	Network.disconnect_message(MessageTypes.GAME_END)
	Network.disconnect_message(MessageTypes.PLAYER_DISCONNECTED)

	get_tree().get_root().find_node("InGame", true, false).queue_free()
	Menu.reset_all()
	Menu.get(Menu.MENU_LEVEL.END_GAME, false).score_summary = data[0].response.data.score_summary
	Menu.go(Menu.MENU_LEVEL.END_GAME)