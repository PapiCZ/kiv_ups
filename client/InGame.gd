extends Node2D

onready var Spaceship = preload("res://objects/Spaceship.tscn")
onready var Projectile = preload("res://objects/Projectile.tscn")
onready var Asteroid = preload("res://objects/Asteroid.tscn")
onready var Score = preload("res://objects/Score.tscn")

var last_nodes = {}
var active_nodes = {}
var score_offset = 0
var player_scores = {}

func _ready():
	Network.connect("disconnected", self, "_network_disconnect")

func update_game_tree(parent_obj, node):
	var game_node = null

	if parent_obj.has_node(node.id):
		# update
		game_node = parent_obj.get_node(node.id)
		if node.type == "spaceship":
			update_spaceship(game_node, node)
		elif node.type == "projectile":
			update_projectile(game_node, node)
		elif node.type == "asteroid":
			update_asteroid(game_node, node)
		elif node.type == "score":
			update_score(game_node, node)
	else:
		# create
		print("Creating new node: ", node.id)

		if node.type == "spaceship":
			game_node = create_spaceship(node)
			update_spaceship(game_node, node)
		elif node.type == "projectile":
			game_node = create_projectile(node)
			update_projectile(game_node, node)
		elif node.type == "asteroid":
			game_node = create_asteroid(node)
			update_asteroid(game_node, node)
		elif node.type == "score":
			game_node = create_score(node)
			update_score(game_node, node)
		else:
			game_node = Node2D.new()

		if game_node != null:
			game_node.set_name(node.id)
			parent_obj.add_child(game_node)

	active_nodes[node.id] = [parent_obj, node.id]

	if node.children != null:
		for child_node in node.children:
			update_game_tree(game_node, child_node)
		

func create_spaceship(node):
	var spaceship = Spaceship.instance()
	spaceship.player_name = node.value.player_name

	return spaceship

func update_spaceship(spaceship, node):
	if node.value.player_name == Network.username:
		return

	spaceship.position = Vector2(node.value.pos_x, node.value.pos_y)
	spaceship.velocity = Vector2(node.value.velocity_x, node.value.velocity_y)
	spaceship.rotation = node.value.rotation

func create_projectile(node):
	var projectile = Projectile.instance()

	return projectile

func update_projectile(projectile, node):
	projectile.position = Vector2(node.value.pos_x, node.value.pos_y)
	projectile.velocity = Vector2(node.value.velocity_x, node.value.velocity_y)
	projectile.rotation = node.value.rotation

func create_asteroid(node):
	var asteroid = Asteroid.instance()
	asteroid.scale = Vector2(node.value.scale, node.value.scale)

	return asteroid

func update_asteroid(asteroid, node):
	asteroid.position = Vector2(node.value.pos_x, node.value.pos_y)
	asteroid.velocity = Vector2(node.value.velocity_x, node.value.velocity_y)
	asteroid.rotation = node.value.rotation

func create_score(node):
	var score = Score.instance()
	score.position = Vector2(get_viewport_rect().size.x - 220, 50 + score_offset)
	score_offset += 40
	score.update_player_name(node.value.player_name)
	player_scores[node.value.player_name] = score

	return score

func update_score(score, node):
	score.update_score(node.value.score)

func _update_state(data):
	update_game_tree(self, data[0].response.data.game_tree)

	# Remove active_nodes from last_nodes
	for key in active_nodes.keys():
		last_nodes.erase(key)
	
	for value in last_nodes.values():
		value[0].get_node(value[1]).queue_free()

	last_nodes = active_nodes
	active_nodes = {}

func _player_disconnected(data):
	var player_name = data[0].response.data.player_name

	if player_scores.has(player_name) and player_scores[player_name] != null:
		player_scores[player_name].set_disconnected()

func _network_disconnect():
	Network.disconnect_message(MessageTypes.UPDATE_STATE)
	Network.disconnect_message(MessageTypes.GAME_END)
	Network.disconnect_message(MessageTypes.PLAYER_DISCONNECTED)
	queue_free()
	Menu.go(Menu.MENU_LEVEL.MAIN)