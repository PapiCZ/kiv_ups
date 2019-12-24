extends KinematicBody2D

export var max_speed = 300
export var max_rotation_speed = 150
var velocity = Vector2()
var player_name

func _init():
	self.player_name = player_name

	Network.connect_message(MessageTypes.UPDATE_STATE, self, "_update_state")

func _physics_process(delta):
	var original_position = position
	var original_velocity = velocity
	var original_rotation = rotation

	if Network.username != player_name:
		return

	if Input.is_action_pressed("ui_up"):
		velocity.x = 0
		velocity.y = -1
		velocity = velocity.rotated(rotation)
	elif Input.is_action_pressed("ui_down"):
		velocity.x = 0
		velocity.y = 1
		velocity = velocity.rotated(rotation)
	else:
		velocity.x = 0
		velocity.y = 0

	if Input.is_action_pressed("ui_left"):
		rotation -= deg2rad(max_rotation_speed) * delta
		velocity = velocity.rotated(-deg2rad(max_rotation_speed) * delta)
	if Input.is_action_pressed("ui_right"):
		rotation += deg2rad(max_rotation_speed) * delta
		velocity = velocity.rotated(deg2rad(max_rotation_speed) * delta)

	velocity = velocity * max_speed * delta	
	move_and_collide(velocity)

	if position != original_position or velocity != original_velocity or rotation != original_rotation: 
		Network.send({
			"pos_x": int(position.x),
			"pos_y": int(position.y),
			"velocity_x": velocity.x,
			"velocity_y": velocity.y,
			"rotation": rotation,
		}, MessageTypes.PLAYER_MOVE)

func _update_state(data):
	for spaceship_node in data[0].response.data.game_tree.children:
		if spaceship_node.type != "spaceship":
			continue

		var spaceship = spaceship_node.value
		if spaceship.player_name == player_name and Network.username != player_name:
			position.x = spaceship.pos_x
			position.y = spaceship.pos_y
			rotation = spaceship.rotation