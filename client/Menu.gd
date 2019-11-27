extends Node

enum MENU_LEVEL {
	LOGIN,
	MAIN,
	START_GAME,
	CREATE_LOBBY,
	WAITING_LOBBY,
}

var menus = {
	MENU_LEVEL.LOGIN : preload("res://forms/login_menu/Login.tscn").instance(),
	MENU_LEVEL.MAIN : preload("res://forms/main_menu/MainMenu.tscn").instance(),
	MENU_LEVEL.START_GAME : preload("res://forms/start_game_menu/StartGame.tscn").instance(),
	MENU_LEVEL.CREATE_LOBBY : preload("res://forms/create_lobby_menu/CreateLobby.tscn").instance(),
	MENU_LEVEL.WAITING_LOBBY : preload("res://forms/waiting_lobby_menu/WaitingLobby.tscn").instance(),
}

func _ready():
	pass

func all():
	return menus

func get(menu):
	return menus[menu]

func reset(menu):
	menus[menu] = load(get(menu).filename).instance()