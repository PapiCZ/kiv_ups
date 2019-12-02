extends VBoxContainer

var lobby_name

# Called when the node enters the scene tree for the first time.
func _ready():
	pass # Replace with function body.

# Called every frame. 'delta' is the elapsed time since the previous frame.
#func _process(delta):
#	pass

func _on_Back_pressed():
	Network.send({
		"name": lobby_name
	}, MessageTypes.DELETE_LOBBY)
	Menu.back()
