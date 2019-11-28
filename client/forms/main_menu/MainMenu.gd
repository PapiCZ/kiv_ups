extends VBoxContainer

# Called when the node enters the scene tree for the first time.
func _ready():
# warning-ignore:return_value_discarded
	$FormContainer/StartGame.connect("pressed", self, "_on_StartGame_pressed")
# warning-ignore:return_value_discarded
	$FormContainer/Leaderboards.connect("pressed", self, "_on_Leaderboards_pressed")
# warning-ignore:return_value_discarded
	$FormContainer/SignOut.connect("pressed", self, "_on_SignOut_pressed")
# warning-ignore:return_value_discarded
	$FormContainer/Quit.connect("pressed", self, "_on_Quit_pressed")

func _on_StartGame_pressed():
	Menu.go(Menu.MENU_LEVEL.START_GAME)

func _on_SignOut_pressed():
	Menu.back()
	
func _on_Quit_pressed():
	Network.stop()
	get_tree().quit()

# Called every frame. "delta" is the elapsed time since the previous frame.
#func _process(delta):
#	pass
