extends Node2D

const PARTICLES_SPEED_SCALE_IDLE = 1
const PARTICLES_SPEED_SCALE_MOVE = 3

export var friction = 0.7
export var max_speed = 500
export var rotation_speed = 200

var speed = 0
var velocity = Vector2()
var player_name
var screen_size = null
var can_shoot = true

onready var Projectile = preload("res://objects/Projectile.tscn")

onready var nlast_position = position
onready var nlast_velocity = velocity
onready var nlast_rotation = rotation

func _init():
	self.player_name = player_name

func _ready():
	screen_size = get_viewport_rect().size
	$PlayerName/Label.text = self.player_name
	$ProjectileTimer.connect("timeout", self, "_on_ProjectileTimer_timeout")

func _process(delta):
	$PlayerName.rotation = -rotation

	var particles_speed_scale = PARTICLES_SPEED_SCALE_IDLE

	if position != nlast_position:
		particles_speed_scale = PARTICLES_SPEED_SCALE_MOVE

	$Particles2D.speed_scale = particles_speed_scale

	if Network.username != player_name:
		return

	velocity.x = 0
	velocity.y = -1
	velocity = velocity.rotated(rotation)

	if Input.is_action_pressed("ui_up"):
		speed += max_speed * friction * delta
	else:
		# Breaking
		speed -= max_speed * friction * delta

	if speed > max_speed:
		speed = max_speed
	elif speed < 0:
		speed = 0

	if Input.is_action_pressed("ui_left"):
		rotation -= deg2rad(rotation_speed) * delta
		velocity = velocity.rotated(-deg2rad(rotation_speed) * delta)
	if Input.is_action_pressed("ui_right"):
		rotation += deg2rad(rotation_speed) * delta
		velocity = velocity.rotated(deg2rad(rotation_speed) * delta)

	velocity *= speed
	position += velocity * delta
	if clamp(position.x, 0, screen_size.x) != position.x \
		or clamp(position.y, 0, screen_size.y) != position.y:
		
		# Collision
		position -= velocity * delta
		speed = 0

	if Input.is_action_pressed("ui_shoot") and can_shoot:
		$ProjectileTimer.start()
		can_shoot = false

		Network.send({}, MessageTypes.SHOOT_PROJECTILE)

func _physics_process(delta):

	if (position != nlast_position or velocity != nlast_velocity or rotation != nlast_rotation) \
		and Network.username == player_name:
		Network.send({
			"pos_x": int(position.x),
			"pos_y": int(position.y),
			"velocity_x": velocity.x,
			"velocity_y": velocity.y,
			"rotation": rotation,
		}, MessageTypes.PLAYER_MOVE)

	nlast_position = position
	nlast_velocity = velocity
	nlast_rotation = rotation

func _on_ProjectileTimer_timeout():
	$ProjectileTimer.stop()
	can_shoot = true
