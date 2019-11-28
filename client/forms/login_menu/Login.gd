extends VBoxContainer

# Called when the node enters the scene tree for the first time.
func _ready():
# warning-ignore:return_value_discarded
	$FormContainer/LoginButton.connect("pressed", self, "_on_LoginButton_pressed")
# warning-ignore:return_value_discarded
	$FormContainer/Quit.connect("pressed", self, "_on_Quit_pressed")

func _on_LoginButton_pressed():
	if len($FormContainer/Nickname.text):
		Menu.go(Menu.MENU_LEVEL.MAIN)
		Network.set_auth_data($FormContainer/Nickname.text)
		Network.start_thread("127.0.0.1", 35000)

func _on_Quit_pressed():
	get_tree().quit()