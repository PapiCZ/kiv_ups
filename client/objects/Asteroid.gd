extends Node2D

#warning-ignore-all:unused_class_variable

var velocity
var rotation_

func _process(delta):
	if velocity != null:
		position += velocity * delta