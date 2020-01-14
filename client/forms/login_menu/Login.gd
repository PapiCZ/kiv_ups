extends VBoxContainer

func _ready():
	Network.connect("authenticated", self, "_authenticated")
	Network.connect("authentication_failed", self, "_authentication_failed")

func _on_LoginButton_pressed():
	if len($FormContainer/Nickname.text):
		Network.stop(true)
		Network.set_auth_data($FormContainer/Nickname.text)
		Network.start_thread($FormContainer/Host.text, $FormContainer/Port.value)

func _authenticated():
	Menu.go(Menu.MENU_LEVEL.MAIN)

func _authentication_failed():
	var authentication_failed_dialog = get_tree().get_root().get_node("Game/AuthenticationFailedDialog")
	authentication_failed_dialog.popup_centered()
	$FormContainer/Nickname.text = ""
	Network.stop(true)

func _on_Quit_pressed():
	get_tree().quit()