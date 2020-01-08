extends Node

signal change_menu

enum MENU_LEVEL {
	LOGIN,
	MAIN,
	START_GAME,
	CREATE_LOBBY,
	LOBBY,
	JOIN_LOBBY,
	END_GAME,
}

var menus = {
	MENU_LEVEL.LOGIN : preload("res://forms/login_menu/Login.tscn").instance(),
	MENU_LEVEL.MAIN : preload("res://forms/main_menu/MainMenu.tscn").instance(),
	MENU_LEVEL.START_GAME : preload("res://forms/start_game_menu/StartGame.tscn").instance(),
	MENU_LEVEL.CREATE_LOBBY : preload("res://forms/create_lobby_menu/CreateLobby.tscn").instance(),
	MENU_LEVEL.LOBBY : preload("res://forms/lobby_menu/Lobby.tscn").instance(),
	MENU_LEVEL.JOIN_LOBBY : preload("res://forms/join_lobby_menu/JoinLobby.tscn").instance(),
	MENU_LEVEL.END_GAME : preload("res://forms/end_game_menu/EndGame.tscn").instance(),
}

var menu_stack = []

func load():
	for menu in Menu.all().values():
		emit_signal("change_menu", menu)

func all():
	return menus

func get(menu, call_load=true):
	var menu_obj = menus[menu]

	if call_load and menu_obj.has_method("_load"):
		menu_obj._load()

	return menu_obj

func reset(menu):
	# Reload given menu from the scene file
	var menu_obj = load(get(menu, false).filename).instance()
	menus[menu] = menu_obj

func reset_all():
	# Reload all menus
	for menu in menus:
		reset(menu)

func go(menu):
	# Change menu
	print(menu)
	var menu_obj = get(menu)
	print(menu_obj)
	menu_stack.append(menu_obj)
	emit_signal("change_menu", menu_obj)

func back():
	# Go back
	if len(menu_stack) >= 2:
		menu_stack.pop_back()
		emit_signal("change_menu", menu_stack.back())

func hide(menu):
	menu.visible = false

func hide_current():
	hide(menu_stack[len(menu_stack) - 1])

func show(menu):
	menu.visible = true

func show_current():
	# Hide menu that is on top of the menu_stack
	show(menu_stack[len(menu_stack) - 1])

func hide_and_reset_stack():
	menu_stack = []
	emit_signal("change_menu", null)
