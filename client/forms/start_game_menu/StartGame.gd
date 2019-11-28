extends VBoxContainer

# Called when the node enters the scene tree for the first time.
func _ready():
# warning-ignore:return_value_discarded
	$FormContainer/CreateLobby.connect("pressed", self, "_on_CreateLobby_pressed")
# warning-ignore:return_value_discarded
	$FormContainer/JoinLobby.connect("pressed", self, "_on_JoinLobby_pressed")
# warning-ignore:return_value_discarded
	$FormContainer/Back.connect("pressed", self, "_on_Back_pressed")

func _on_CreateLobby_pressed():
	Menu.go(Menu.MENU_LEVEL.CREATE_LOBBY)

func _on_JoinLobby_pressed():
	Menu.go(Menu.MENU_LEVEL.JOIN_LOBBY)

func _on_Back_pressed():
	Menu.back()