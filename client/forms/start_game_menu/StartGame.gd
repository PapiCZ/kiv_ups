extends VBoxContainer

signal change_menu

# Called when the node enters the scene tree for the first time.
func _ready():
# warning-ignore:return_value_discarded
	$FormContainer/CreateLobby.connect("pressed", self, "_on_CreateLobby_pressed")
# warning-ignore:return_value_discarded
	$FormContainer/JoinLobby.connect("pressed", self, "_on_JoinLobby_pressed")
# warning-ignore:return_value_discarded
	$FormContainer/Back.connect("pressed", self, "_on_Back_pressed")

func _on_CreateLobby_pressed():
	emit_signal("change_menu", load("res://forms/create_lobby_menu/CreateLobby.tscn").instance())

func _on_Back_pressed():
	emit_signal("change_menu", load("res://forms/main_menu/MainMenu.tscn").instance())