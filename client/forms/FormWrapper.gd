extends CenterContainer

# Declare member variables here. Examples:
# var a = 2
# var b = "text"

# Called when the node enters the scene tree for the first time.
func _ready():
	Menu.connect("menu_added", self, "_menu_added")
	Menu.load()
	_on_menu_changed(Menu.get(Menu.MENU_LEVEL.LOGIN))

func _menu_added(menu):
	menu.connect("change_menu", self, "_on_menu_changed")

func _on_menu_changed(menu_root_node):
	for i in range(0, get_child_count()):
		call_deferred("remove_child", get_child(i))
	call_deferred("add_child", menu_root_node)
