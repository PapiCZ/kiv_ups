extends VBoxContainer

# Called when the node enters the scene tree for the first time.
func _ready():
# warning-ignore:return_value_discarded
	$FormContainer/CreateLobby.connect("pressed", self, "_on_CreateLobby_pressed")
# warning-ignore:return_value_discarded
	$FormContainer/Back.connect("pressed", self, "_on_Back_pressed")

func _on_Back_pressed():
	Menu.back()
	Menu.reset(Menu.MENU_LEVEL.CREATE_LOBBY)

func _on_CreateLobby_pressed():
	var lobby_name = $FormContainer/LobbyName.text

	if len(lobby_name):
		Network.send({
			"name": lobby_name,
			"players_limit": $FormContainer/PlayersLimit.value,
		}, MessageTypes.CREATE_LOBBY, self, "_on_lobby_created", [lobby_name])

func _on_lobby_created(data):
	var menu = Menu.get(Menu.MENU_LEVEL.LOBBY, false)
	menu.lobby_name = data[1]
	menu.get_node("Label").text = "Lobby " + data[1]
	Menu.go(Menu.MENU_LEVEL.LOBBY)