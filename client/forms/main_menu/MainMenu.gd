extends VBoxContainer

var InGame = preload("res://InGame.tscn")

func _load():
	var t = Timer.new()
	t.set_wait_time(0.5)
	add_child(t)
	t.start()
	yield(t, "timeout")
	Network.send({}, MessageTypes.GAME_RECONNECT_AVAILABLE, self, "_game_reconnect_available")

func _on_StartGame_pressed():
	Menu.go(Menu.MENU_LEVEL.START_GAME)

func _on_SignOut_pressed():
	Menu.back()
	
func _on_Quit_pressed():
	Network.stop()
	get_tree().quit()

func _game_reconnect_available(data):
	print(data[0].response)
	if data[0].response.status and data[0].response.data.available:
		$FormContainer/ReconnectGame.visible = true
		$FormContainer/StartGame.visible = false

func _on_ReconnectGame_pressed():
	Menu.hide_current()
	
	var ingame = InGame.instance()
	ingame.set_name("InGame")
	Network.connect_message(MessageTypes.UPDATE_STATE, ingame, "_update_state")
	Network.connect_message(MessageTypes.GAME_END, self, "_end_game")
	Network.connect_message(MessageTypes.PLAYER_DISCONNECTED, ingame, "_player_disconnected")
	get_parent().get_parent().get_parent().add_child(ingame)

func _end_game(data):
	get_tree().get_root().find_node("InGame", true, false).queue_free()
	Menu.hide_and_reset_stack()
	Menu.go(Menu.MENU_LEVEL.END_GAME)