extends Node

signal disconnected
signal connected
signal authenticated
signal authentication_failed

const DEBUG = true

var host = ""
var port = 0

var username = null

var client = null
var mutex = null
var thread = null
var authenticated = false
var _kill_thread = false
var pending_requests = {}

var KeepAliveTimer = null
const KEEP_ALIVE_INVERVAL = 1

var disconnection_notified = false

onready var PendingRequest = preload("res://networking/PendingRequest.gd")
onready var GameProtocol = preload("res://networking/GameProtocol.gd").new()
onready var ProtoMessage = preload("res://networking/ProtoMessage.gd")
onready var KeepAlive = preload("res://networking/KeepAlive.gd").new()
onready var Utils = preload("res://networking/Utils.gd")

onready var status_sprite = get_tree().get_root().get_node("Game/NetworkStatus/Sprite")
onready var network_debug = get_tree().get_root().get_node("Game/NetworkDebug")
onready var network_ping = get_tree().get_root().get_node("Game/NetworkPing")

var connected_signals = {}

func _ready():
	# Intialize mutex for shared resources
	mutex = Mutex.new()

func connect_message(type, callback_obj, callback_func):
	# Connect given callback to specific message type.
	# Callback is called when given message type comes
	# the server.
	mutex.lock()
	if not connected_signals.has(type):
		connected_signals[type] = []

	connected_signals[type].append([callback_obj, callback_func])
	mutex.unlock()

func disconnect_message(type):
	# Disconnects all callbacks from given message type
	mutex.lock()
	connected_signals[type] = []
	mutex.unlock()

func network_ok():
	# Set network indicator to OK
	if DEBUG:
		network_debug.log("Connected to %s:%d" % [host, port])
	status_sprite.call_deferred("set_texture", preload("res://sprites/network_ok.png"))
	
func network_error():
	# Set network indicator to ERROR
	if DEBUG:
		network_debug.log("Disconnected from %s:%d" % [host, port])
	status_sprite.call_deferred("set_texture", preload("res://sprites/network_error.png"))

func send(message, type, \
	response_callback_obj=null, response_callback_func=null, response_args=null, \
	timeout_callback_obj=null, timeout_callback_func=null, timeout_args=null):
	# Send message to the server

	var proto_message = ProtoMessage.new()
	proto_message.message = message
	proto_message.type = type
	proto_message.request_id = Utils.random_request_id()

	var data = GameProtocol.encode(proto_message)
	if DEBUG:
		network_debug.log("Sent %s" % data.get_string_from_utf8())

	if response_callback_obj != null or timeout_callback_obj != null:
		# Create pending request if user defined any callback
		var pr = PendingRequest.new()
		pr.id = proto_message.request_id
		pr.request = proto_message
		pr.request_time = OS.get_ticks_msec()
		pr._response_callback_obj = response_callback_obj
		pr._response_callback_func = response_callback_func
		pr._response_args = response_args
		pr._timeout_callback_obj = timeout_callback_obj
		pr._timeout_callback_func = timeout_callback_func
		pr._timeout_args = timeout_args

		pr._timeout_timer = Timer.new()
		add_child(pr._timeout_timer)
		pr._timeout_timer.set_wait_time(pr._timeout)
		var _timeout_args = [pr]
		if pr._timeout_args != null:
			for arg in pr._timeout_args:
				_timeout_args.append(arg)
		if pr._timeout_callback_obj != null and pr._timeout_callback_func != null:
			pr._timeout_timer.connect("timeout", pr._timeout_callback_obj, pr._timeout_callback_func, _timeout_args)
		pr._timeout_timer.connect("timeout", self, "_on_RequestTimeoutTimer_timeout", [pr])
		pr._timeout_timer.start()

		mutex.lock()
		if client.get_status() != client.STATUS_CONNECTED:
			mutex.unlock()
			# Check if client is still connected to the server
			# Because we can lose connection between keep-alive messages
			pr.request.queue_free()
			pr._timeout_timer.stop()
			pr._timeout_timer.queue_free()
			if pr._timeout_callback_obj != null and pr._timeout_callback_func != null:
				pr._timeout_callback_obj.call_deferred(pr._timeout_callback_func, _timeout_args)
			pr.call_deferred("free")
			return

		# Send prepared data to the server and add request to pending requests
		pending_requests[pr.id] = pr
		client.put_data(data)
		mutex.unlock()
	else:
		# User didn't define any callbacks. We don't need proto_message anymore.
		proto_message.queue_free()

		mutex.lock()
		if client.get_status() != client.STATUS_CONNECTED:
			mutex.unlock()
			return

		client.put_data(data)
		mutex.unlock()
	
	return data

func _on_RequestTimeoutTimer_timeout(pr):
	# Stop and remove timeout timer from scene tree
	pr._timeout_timer.stop()
	remove_child(pr._timeout_timer)

	# Remove pending request
	mutex.lock()
	pending_requests.erase(pr.id)
	mutex.unlock()

func recv_message_loop():
	# Prepare buffer for incoming messages
	var buff = PoolByteArray()

	while true:
		mutex.lock()
		if _kill_thread:
			mutex.unlock()
			return

		var available_bytes = 0
		if client.get_status() == client.STATUS_CONNECTED:
			available_bytes = client.get_available_bytes()

		var buff_len = len(buff)
		if available_bytes or buff_len:
			# Someone sent us data or we have some unproccessed data
			if available_bytes:
				# We have some new data
				if buff_len + available_bytes > pow(2, 15):
					# Maximal buffer size si 2^15 (is't godot limitation)
					available_bytes = pow(2, 15) - buff_len

				var data_result = client.get_data(available_bytes)
				mutex.unlock()
				buff.append_array(data_result[1])
			else:
				mutex.unlock()

			# Try to decode message
			var result = GameProtocol.decode(buff)

			if result == null or len(result) != 2 or result[0] == null or result[1] == null:
				continue

			var message_len = result[0]
			var proto_message = result[1]

			if DEBUG:
				network_debug.log("Received: %d %s %s" % [proto_message.type, proto_message.request_id, proto_message.message])

			# Process message... call response callbacks and destroy timeout timer
			mutex.lock()
			if pending_requests.has(proto_message.request_id):
				var pr = pending_requests[proto_message.request_id]
				pending_requests.erase(proto_message.request_id)
				mutex.unlock()
				pr.response = proto_message.message
				pr.response_time = OS.get_ticks_msec()
				network_ping.set_text("%d ms" % (pr.response_time - pr.request_time))
				if pr._timeout_timer != null:
					pr._timeout_timer.stop()
					remove_child(pr._timeout_timer)

				if pr._response_callback_obj != null and pr._response_callback_func != null:
					var response_args = [pr]
					if pr._response_args != null:
						for arg in pr._response_args:
							response_args.append(arg)
					pr._response_callback_obj.call_deferred(pr._response_callback_func, response_args)
				pr.request.queue_free()
				pr._timeout_timer.queue_free()
				pr.queue_free()
			else:
				mutex.unlock()

			mutex.lock()

			# Call callbacks that are connected to given type of message
			if connected_signals.has(proto_message.type):
				for callback in connected_signals[proto_message.type]:
					var pr = PendingRequest.new()
					pr.response = proto_message.message
					pr.response_time = OS.get_ticks_msec()
					var response_args = [pr]
					callback[0].call_deferred(callback[1], response_args)
					pr.call_deferred("free")

			proto_message.queue_free()
			mutex.unlock()

			# Remove proccessed message from buffer
			if len(buff) > message_len:
				buff = buff.subarray(message_len, len(buff) - 1)
			else:
				buff = PoolByteArray()
		else:
			mutex.unlock()

func set_auth_data(username):
	# Setter for authentication data. Currently only username.
	self.username = username

func auth(response_callback_obj=null, response_callback_func=null):
	# Send authentication request to the server
	send({"name": username}, MessageTypes.AUTHENTICATE, response_callback_obj, response_callback_func)

func start_thread(host, port, auto_reconnect=true):
	print("Spawned new network thread")
	if init(host, port) == true:
		# Try to authenticate
		auth(self, "_auth_callback")

		# Intialization of connection was successful
		if disconnection_notified:
			get_tree().get_root().get_node("Game/NetworkConnectedDialog").popup_centered()

		disconnection_notified = false
		get_tree().get_root().get_node("Game/NetworkDisconnectedDialog").hide()
		thread = Thread.new()
		thread.start(self, "_start")
	else:
		if auto_reconnect:
			# Sleep 1 second and try to reconnect
			yield(get_tree().create_timer(1), "timeout")
			self.call_deferred("reconnect")

func init(host, port):
	# Initialize server connection and try to connect to the server
	mutex.lock()
	self.host = host
	self.port = port

	client = StreamPeerTCP.new()
	client.connect_to_host(self.host, self.port)
	authenticated = false
	mutex.unlock()
	print("Connecting to host")

	for i in range(10000000):
		if client.get_status() == client.STATUS_CONNECTED:
			break

	if client.get_status() != client.STATUS_CONNECTED:
		return false
	
	network_ok()
	emit_signal("connected")

	return true

func _auth_callback(data):
	if data[0].response.status:
		mutex.lock()
		authenticated = true
		mutex.unlock()
		emit_signal("authenticated")
	else:
		mutex.lock()
		authenticated = false
		mutex.unlock()
		emit_signal("authentication_failed")

func _start(params):
	# Start method that is meant to be started in new thread, because
	# function recv_message_loop is blocking
	mutex.lock()
	KeepAlive.failed = false
	KeepAliveTimer = Timer.new()
	KeepAliveTimer.set_wait_time(KEEP_ALIVE_INVERVAL)
	KeepAliveTimer.connect("timeout", self, "_on_KeepAliveTimer_timeout")
	KeepAliveTimer.start()
	get_tree().get_root().add_child(KeepAliveTimer)
	mutex.unlock()

	recv_message_loop()
	
func stop(ignore_alert=false):
	# Show alert dialog
	var network_disconnected_dialog = get_tree().get_root().get_node("Game/NetworkDisconnectedDialog")
	if not network_disconnected_dialog.visible and not disconnection_notified and not ignore_alert:
		get_tree().get_root().get_node("Game/NetworkConnectedDialog").hide()
		disconnection_notified = true
		network_disconnected_dialog.popup_centered()

	# Kill network thread and disconnect from host
	if thread:
		_kill_thread = true
		thread.wait_to_finish()
		_kill_thread = false

	if client:
		client.disconnect_from_host()

	emit_signal("disconnected")
	network_error()

	# Flush all pending requests
	for pr in pending_requests.values():
		pr._timeout_timer.stop()
		pr._timeout_timer.queue_free()
		pr.request.queue_free()
		pr.queue_free()

	mutex = Mutex.new()
	pending_requests.clear()

	# Kill keep-alive timer
	if KeepAliveTimer != null:
		KeepAliveTimer.stop()
		KeepAliveTimer.queue_free()
		KeepAliveTimer = null

func reconnect():
	print("Reconnecting")
	stop()
	print("Killed network thread")
	self.start_thread(host, port)

func _on_KeepAliveTimer_timeout():
	print("Sending Keep-Alive")
	mutex.lock()
	KeepAlive.send()
	mutex.unlock()