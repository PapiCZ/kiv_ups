extends Node2D

export var speed = 600

var velocity = Vector2(0, -1)

func _process(delta):
	if velocity != null:
		# Interpolate position between server messages
		position += velocity * delta
