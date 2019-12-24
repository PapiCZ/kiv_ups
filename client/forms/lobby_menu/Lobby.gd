extends VBoxContainer

var lobby_name
var InGame = preload("res://InGame.tscn")

onready var Spaceship = preload("res://objects/SpaceShip.tscn")

func _load():
	Network.send({}, MessageTypes.LIST_LOBBY_PLAYERS, self, "_on_lobby_players_loaded")

	Network.connect_message(MessageTypes.START_GAME_RESPONSE, self, "_start_game")

func _on_lobby_players_loaded(data):
	if data[0].response.status:
		for player_name in data[0].response.data.players:
			$FormContainer/PlayersContainer.add_player(player_name)

func _on_Start_pressed():
	Network.disconnect_message(MessageTypes.START_GAME_RESPONSE)
	Network.send({}, MessageTypes.START_GAME, self, "_start_game")

func _start_game(data):
	Menu.hide_and_reset_stack()

	var ingame = InGame.instance()
	var i = 0
	for player in $FormContainer/PlayersContainer.players:
		i += 1
		var spaceship = Spaceship.instance()
		spaceship.player_name = player.name
		spaceship.position.x = 100 * i
		spaceship.position.y = 500
		ingame.add_child(spaceship)

	get_parent().get_parent().get_parent().add_child(ingame)

func _on_Back_pressed():
	Network.send({
		"name": lobby_name
	}, MessageTypes.DELETE_LOBBY)
	Menu.back()
