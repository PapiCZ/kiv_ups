extends Node

signal change_menu

enum MENU_LEVEL {
	LOGIN,
	MAIN,
	START_GAME,
	CREATE_LOBBY,
	WAITING_LOBBY,
	JOIN_LOBBY
}

var menus = {
	MENU_LEVEL.LOGIN : preload("res://forms/login_menu/Login.tscn").instance(),
	MENU_LEVEL.MAIN : preload("res://forms/main_menu/MainMenu.tscn").instance(),
	MENU_LEVEL.START_GAME : preload("res://forms/start_game_menu/StartGame.tscn").instance(),
	MENU_LEVEL.CREATE_LOBBY : preload("res://forms/create_lobby_menu/CreateLobby.tscn").instance(),
	MENU_LEVEL.WAITING_LOBBY : preload("res://forms/waiting_lobby_menu/WaitingLobby.tscn").instance(),
	MENU_LEVEL.JOIN_LOBBY : preload("res://forms/join_lobby_menu/JoinLobby.tscn").instance(),
}

var menu_stack = []

func load():
	for menu in Menu.all().values():
		emit_signal("menu_added", menu)

func all():
	return menus

func get(menu, call_load=true):
	var menu_obj = menus[menu]

	if call_load and menu_obj.has_method("_load"):
		menu_obj._load()

	return menu_obj

func reset(menu):
	var menu_obj = load(get(menu, false).filename).instance()
	menus[menu] = menu_obj
	emit_signal("menu_added", menu_obj)

func go(menu):
	var menu_obj = get(menu)
	menu_stack.append(menu_obj)
	emit_signal("change_menu", menu_obj)

func back():
	if len(menu_stack) >= 2:
		menu_stack.pop_back()
		emit_signal("change_menu", menu_stack.back())
