extends RichTextLabel

func log(message):
	# Ignore keepalive
	if "|100|15" in message or "Received: 100" in message:
		return

	call_deferred("_log", message)

func _log(message):
	text += str(message) + "\n"


