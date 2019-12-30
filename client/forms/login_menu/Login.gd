extends VBoxContainer

func _on_LoginButton_pressed():
	if len($FormContainer/Nickname.text):
		Network.set_auth_data($FormContainer/Nickname.text)
		Network.start_thread($FormContainer/Host.text, $FormContainer/Port.value)
		Menu.go(Menu.MENU_LEVEL.MAIN)

func _on_Quit_pressed():
	get_tree().quit()