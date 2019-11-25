extends RichTextLabel

func log(message):
	call_deferred("_log", message)

func _log(message):
	text += str(message) + "\n"


