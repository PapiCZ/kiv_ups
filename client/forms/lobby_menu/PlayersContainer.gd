extends HBoxContainer

var players = []

func add_player(player_name):	
	var player_container = VBoxContainer.new()
	
	var icon = TextureRect.new()
	icon.texture = load("res://sprites/player_not_ready.png")

	var player_name_label = Label.new()
	player_name_label.text = player_name
	player_name_label.align = player_name_label.ALIGN_CENTER

	player_container.add_child(icon)
	player_container.add_child(player_name_label)

	add_child(player_container)

	players.append({
		"name": player_name,
		"node": player_container
	})

func remove_player(player_name):
	for player in players:
		if player_name == player.name:
			players.erase(player)
			player.node.queue_free()
			break