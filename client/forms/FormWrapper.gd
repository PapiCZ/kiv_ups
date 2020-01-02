extends CenterContainer

func _ready():
	Menu.connect("change_menu", self, "_on_menu_changed")
	_on_menu_changed(Menu.get(Menu.MENU_LEVEL.LOGIN))

func _on_menu_changed(menu_root_node):
	for i in range(0, get_child_count()):
		call_deferred("remove_child", get_child(i))

	if menu_root_node:
		call_deferred("add_child", menu_root_node)
