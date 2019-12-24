extends Node

func send():
	Network.send({
		"ping": "pong"
	}, MessageTypes.KEEP_ALIVE, self, "_check", null, self, "_network_error", null)

func _check(data):
	var response = data[0].response
	if not (typeof(response) == TYPE_DICTIONARY and response.has("data") \
		and response["data"].has("ping") and response["data"]["ping"] == "ping-pong"):
		Network.network_error()

func _network_error(data):
	Network.network_error()
	Network.reconnect()