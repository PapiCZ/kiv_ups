extends CenterContainer

# Declare member variables here. Examples:
# var a = 2
# var b = "text"

# Called when the node enters the scene tree for the first time.
func _ready():
	for menu in Menu.all().values():
		menu.connect("change_menu", self, "_on_menu_changed")

	_on_menu_changed(Menu.get(Menu.MENU_LEVEL.LOGIN))

func _on_menu_changed(menu_root_node):
	menu_root_node.connect("change_menu", self, "_on_menu_changed")
	for i in range(0, get_child_count()):
		remove_child(get_child(i))
	add_child(menu_root_node)
