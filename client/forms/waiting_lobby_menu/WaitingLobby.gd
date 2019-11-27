extends VBoxContainer

signal change_menu

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
	}, MessageTypes.DELETE_LOBBY, self, "_on_lobby_deleted")

func _on_lobby_deleted(data):
	emit_signal("change_menu", Menu.get(Menu.MENU_LEVEL.CREATE_LOBBY))
