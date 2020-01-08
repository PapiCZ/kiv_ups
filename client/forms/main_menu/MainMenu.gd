extends VBoxContainer

var InGame = preload("res://InGame.tscn")

func _load():
	if Network.authenticated:
		Network.send({}, MessageTypes.GAME_RECONNECT_AVAILABLE, self, "_game_reconnect_available")
	else:
		Network.connect("authenticated", self, "_network_authenticated")

func _network_authenticated():
	Network.send({}, MessageTypes.GAME_RECONNECT_AVAILABLE, self, "_game_reconnect_available")

func _on_StartGame_pressed():
	Menu.go(Menu.MENU_LEVEL.START_GAME)

func _on_SignOut_pressed():
	Network.stop()
	Menu.go(Menu.MENU_LEVEL.LOGIN)
	
func _on_Quit_pressed():
	Network.stop()
	get_tree().quit()

func _game_reconnect_available(data):
	if data[0].response.status and data[0].response.data.available:
		# Game reconnect is available, let's show or hide some buttons
		$FormContainer/ReconnectGame.visible = true
		$FormContainer/LeaveGame.visible = true
		$FormContainer/StartGame.visible = false

func _on_ReconnectGame_pressed():
	Network.send({}, MessageTypes.RECONNECT)

	Menu.hide_current()
	
	# Start game
	var ingame = InGame.instance()
	ingame.set_name("InGame")
	Network.connect_message(MessageTypes.UPDATE_STATE, ingame, "_update_state")
	Network.connect_message(MessageTypes.GAME_END, self, "_end_game")
	Network.connect_message(MessageTypes.PLAYER_DISCONNECTED, ingame, "_player_disconnected")
	Network.connect_message(MessageTypes.PLAYER_CONNECTED, ingame, "_player_connected")
	get_parent().get_parent().get_parent().add_child(ingame)

func _input(ev):
	if Input.is_action_pressed("ui_leave_ingame") and get_tree().get_root().find_node("InGame", true, false):
		# End the game and show winner
		Network.disconnect_message(MessageTypes.UPDATE_STATE)
		Network.disconnect_message(MessageTypes.GAME_END)
		Network.disconnect_message(MessageTypes.PLAYER_DISCONNECTED)
		Network.disconnect_message(MessageTypes.PLAYER_CONNECTED)

		print(Network.connected_signals)

		get_tree().get_root().find_node("InGame", true, false).queue_free()
		Menu.reset_all()
		Menu.go(Menu.MENU_LEVEL.MAIN)

func _end_game(data):
	# End game and show winner
	Network.disconnect_message(MessageTypes.UPDATE_STATE)
	Network.disconnect_message(MessageTypes.GAME_END)
	Network.disconnect_message(MessageTypes.PLAYER_DISCONNECTED)

	get_tree().get_root().find_node("InGame", true, false).queue_free()
	Menu.reset_all()
	Menu.get(Menu.MENU_LEVEL.END_GAME, false).score_summary = data[0].response.data.score_summary
	Menu.go(Menu.MENU_LEVEL.END_GAME)

func _on_LeaveGame_pressed():
	# Definitely leave the game without possibility of reconnecting
	Network.send({}, MessageTypes.LEAVE_GAME, self, "_leave_game_response")

func _leave_game_response(data):
	if data[0].response.status:
		$FormContainer/ReconnectGame.visible = false
		$FormContainer/LeaveGame.visible = false
		$FormContainer/StartGame.visible = true
