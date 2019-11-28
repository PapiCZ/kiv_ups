extends VBoxContainer

# Called when the node enters the scene tree for the first time.
func _ready():
# warning-ignore:return_value_discarded
	$FormContainer/Back.connect("pressed", self, "_on_Back_pressed")

func _load():
	Network.send({}, MessageTypes.LIST_LOBBIES, self, "_on_lobbies_loaded")

func _on_lobbies_loaded(data):
	
	for lobby in data[0].response.data.lobbies:
		var button = Button.new()
		button.text = "%-29s%d/%d" % [lobby.name, lobby.connected_players, lobby.players_limit]

		$FormContainer/PanelContainer/ScrollContainer/LobbyList.add_child(button)

func _on_Back_pressed():
	Menu.back()
	Menu.reset(Menu.MENU_LEVEL.JOIN_LOBBY)