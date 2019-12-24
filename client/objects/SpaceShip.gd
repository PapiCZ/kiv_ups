extends KinematicBody2D

export var max_speed = 300
export var max_rotation_speed = 150
var velocity = Vector2()

func _process(delta):
	var original_position = position
	var original_velocity = velocity
	var original_rotation = rotation

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

