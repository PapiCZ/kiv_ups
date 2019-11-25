extends VBoxContainer

signal change_menu

# Called when the node enters the scene tree for the first time.
func _ready():
# warning-ignore:return_value_discarded
	$FormContainer/LoginButton.connect("pressed", self, "_on_LoginButton_pressed")
# warning-ignore:return_value_discarded
	$FormContainer/Quit.connect("pressed", self, "_on_Quit_pressed")

func _on_LoginButton_pressed():
	emit_signal("change_menu", load("res://forms/main_menu/MainMenu.tscn").instance())
	Network.set_auth_data($FormContainer/Nickname.text)
	Network.start_thread("127.0.0.1", 35000)

func _on_Quit_pressed():
	get_tree().quit()