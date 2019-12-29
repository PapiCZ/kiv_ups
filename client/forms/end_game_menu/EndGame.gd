extends VBoxContainer

func _on_MainMenu_pressed():
	Network.disconnect_message(MessageTypes.START_GAME_RESPONSE)
	Network.disconnect_message(MessageTypes.LOBBY_PLAYER_CONNECTED)
	Menu.reset_all()
	Menu.go(Menu.MENU_LEVEL.MAIN)
