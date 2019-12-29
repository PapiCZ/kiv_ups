extends VBoxContainer

func _on_CreateLobby_pressed():
	Menu.go(Menu.MENU_LEVEL.CREATE_LOBBY)

func _on_JoinLobby_pressed():
	Menu.go(Menu.MENU_LEVEL.JOIN_LOBBY)

func _on_Back_pressed():
	Menu.back()