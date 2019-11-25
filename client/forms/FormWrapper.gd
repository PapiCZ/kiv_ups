extends CenterContainer

# Declare member variables here. Examples:
# var a = 2
# var b = "text"

# Called when the node enters the scene tree for the first time.
func _ready():
	_on_change_menu(load("res://forms/login_menu/Login.tscn").instance())

func _on_change_menu(menu_root_node):	
	menu_root_node.connect("change_menu", self, "_on_change_menu")
	for i in range(0, get_child_count()):
		get_child(i).queue_free()
	add_child(menu_root_node)
