extends Node2D

var velocity

func _process(delta):
	if velocity != null:
		# Interpolate position between server messages
		position += velocity * delta