extends VBoxContainer

func _ready():
	Network.connect("authenticated", self, "_authenticated")

func _on_LoginButton_pressed():
	if len($FormContainer/Nickname.text):
		Network.stop(true)
		Network.set_auth_data($FormContainer/Nickname.text)
		Network.start_thread($FormContainer/Host.text, $FormContainer/Port.value)

func _authenticated():
	Menu.go(Menu.MENU_LEVEL.MAIN)

func _on_Quit_pressed():
	get_tree().quit()