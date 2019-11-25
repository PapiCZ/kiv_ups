extends VBoxContainer

signal change_menu

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
	emit_signal("change_menu", load("res://forms/start_game_menu/StartGame.tscn").instance())

func _on_SignOut_pressed():
	Network.stop_and_reset()
	emit_signal("change_menu", load("res://forms/login_menu/Login.tscn").instance())
	
func _on_Quit_pressed():
	get_tree().quit()
# Called every frame. 'delta' is the elapsed time since the previous frame.
#func _process(delta):
#	pass


func _on_Timer_timeout():
	Network.send({
		"name": "foobar"
	}, 200, self)