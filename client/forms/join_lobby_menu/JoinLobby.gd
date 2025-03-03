extends VBoxContainer

func _load():
	_on_Refresh_pressed()

func _on_lobbies_loaded(data):
	
	for lobby in data[0].response.data.lobbies:
		var button = Button.new()
		button.text = "%-29s%d/%d" % [lobby.name, lobby.connected_players, lobby.players_limit]
		button.connect("pressed", self, "_connect_lobby", [lobby.name])

		$FormContainer/PanelContainer/ScrollContainer/LobbyList.add_child(button)

func _connect_lobby(lobby_name):
	Menu.reset(Menu.MENU_LEVEL.LOBBY)
	Network.send({
		"name": lobby_name
	}, MessageTypes.JOIN_LOBBY, self, "_on_lobby_connected", [lobby_name])


func _on_lobby_connected(data):
	# Move player to the lobby
	if data[0].response.status:
		var menu = Menu.get(Menu.MENU_LEVEL.LOBBY, false)
		menu.lobby_name = data[1]
		menu.get_node("Label").text = "Lobby " + data[1]
		Menu.go(Menu.MENU_LEVEL.LOBBY)
		Menu.reset(Menu.MENU_LEVEL.LOBBY)

func _on_Back_pressed():
	Menu.back()
	Menu.reset(Menu.MENU_LEVEL.JOIN_LOBBY)

func _on_Refresh_pressed():
	# Refresh list of lobbies
	for i in range(0, $FormContainer/PanelContainer/ScrollContainer/LobbyList.get_child_count()):
		$FormContainer/PanelContainer/ScrollContainer/LobbyList.get_child(i).queue_free()
	Network.send({}, MessageTypes.LIST_LOBBIES, self, "_on_lobbies_loaded")
