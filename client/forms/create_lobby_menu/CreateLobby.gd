extends VBoxContainer

signal change_menu

# Called when the node enters the scene tree for the first time.
func _ready():
# warning-ignore:return_value_discarded
	$FormContainer/CreateLobby.connect("pressed", self, "_on_JoinLobby_pressed")
# warning-ignore:return_value_discarded
	$FormContainer/Back.connect("pressed", self, "_on_Back_pressed")

func _on_Back_pressed():
	emit_signal("change_menu", load("res://forms/start_game_menu/StartGame.tscn").instance())