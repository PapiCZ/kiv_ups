extends Node

var asteroids = []

# Called when the node enters the scene tree for the first time.
func _ready():
	var asteroid = load("res://objects/Asteroid.tscn")
	for _i in range(0, 10):
		var a = asteroid.instance()
		reset_asteroid(a)
		$ParallaxBackground/AsteroidParent.add_child(a)
		asteroids.append(a)
		
func random_asteroid_position(rect):
	var side = randi() % 4
	var position
	
	if side == 0:
		position = Vector2(rand_range(0, rect[0]), -100)
	elif side == 1:
		position = Vector2(rect[0] + 100, rand_range(0, rect[1]))
	elif side == 2:
		position = Vector2(rand_range(0, rect[0]), rect[1] + 100)
	elif side == 3:
		position = Vector2(-100, rand_range(0, rect[1]))
		
	return position
	
func random_asteroid_velocity():
	return Vector2(rand_range(-40, 40), rand_range(-40, 40))
	
func reset_asteroid(asteroid):
	pass
	# var rect = OS.get_screen_size()
	# asteroid.position = random_asteroid_position(rect)
	# asteroid.rotation_ = rand_range(-20, 20)
	# var scale = rand_range(0.5, 1.5)
	# asteroid.get_node("Sprite").scale = Vector2(scale, scale)
	# asteroid.velocity = random_asteroid_velocity()

# Called every frame. "delta" is the elapsed time since the previous frame.
func _process(delta):
	pass
	# for a in asteroids:
		# a.position += a.velocity * delta
		# a.get_node("Sprite").rotation_degrees += a.rotation_ * delta
		
		# var size = OS.get_screen_size()
		# var r = Rect2(-100, -100, size[0] + 200, size[1] + 200)
		
		# if not r.has_point(a.position):
		# 	reset_asteroid(a)