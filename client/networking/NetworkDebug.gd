extends RichTextLabel

func log(message):
	# We dont want to modify this node from network thread
	call_deferred("_log", message)

func _log(message):
	if len(text) > 100000:
		# RichTextLabel can't handle millions of characters
		text = text.right(1000)

	if "asteroidwrapper" in message:
		# Ignore in-game messages. There are tons of them and it costs a lot of performance.
		return

	text += str(message) + "\n"


