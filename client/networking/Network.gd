extends Node

var host = ""
var port = 0

var username = null

var client = null
var mutex = null
var thread = null
var _kill_thread = false
var pending_requests = {}

var KeepAliveTimer = null
const KEEP_ALIVE_INVERVAL = 1

var PendingRequest = preload("res://networking/PendingRequest.gd")
var GameProtocol = preload("res://networking/GameProtocol.gd").new()
var ProtoMessage = preload("res://networking/ProtoMessage.gd")
var KeepAlive = preload("res://networking/KeepAlive.gd").new()
var Utils = preload("res://networking/Utils.gd")
onready var status_sprite = get_tree().get_root().get_node("Game/NetworkStatus/Sprite")
onready var network_debug = get_tree().get_root().get_node("Game/NetworkDebug")
onready var network_ping = get_tree().get_root().get_node("Game/NetworkPing")

func _ready():
	mutex = Mutex.new()

func network_ok():
	mutex.lock()
	network_debug.log("Connected to %s:%d" % [host, port])
	status_sprite.call_deferred("set_texture", preload("res://sprites/network_ok.png"))
	mutex.unlock()
	
func network_error():
	mutex.lock()
	network_debug.log("Disconnected from %s:%d" % [host, port])
	status_sprite.call_deferred("set_texture", preload("res://sprites/network_error.png"))
	mutex.unlock()

func send(message, type, \
	response_callback_obj=null, response_callback_func=null, response_args=null, \
	timeout_callback_obj=null, timeout_callback_func=null, timeout_args=null):
	var proto_message = ProtoMessage.new()
	proto_message.message = message
	proto_message.type = type
	proto_message.request_id = Utils.random_request_id()

	var data = GameProtocol.encode(proto_message)
	network_debug.log("Sent %s" % data.get_string_from_utf8())

	var pr = PendingRequest.new()
	pr.id = proto_message.request_id
	pr.request = proto_message
	pr.request_time = OS.get_ticks_msec()
	pr._response_callback_obj = response_callback_obj
	pr._response_callback_func = response_callback_func
	pr._response_args = [1, 3, 2]
	pr._timeout_callback_obj = timeout_callback_obj
	pr._timeout_callback_func = timeout_callback_func
	pr._timeout_args = timeout_args

	if pr._timeout != null:
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
	pending_requests[pr.id] = pr
	client.put_data(data)
	mutex.unlock()
	
	return data

func _on_RequestTimeoutTimer_timeout(pr):
	pr._timeout_timer.stop()
	remove_child(pr._timeout_timer)

	mutex.lock()
	pending_requests.erase(pr.id)
	mutex.unlock()

func recv_message_loop():
	var buff = PoolByteArray()

	while true:
		mutex.lock()
		if _kill_thread:
			mutex.unlock()
			return

		var available_bytes = client.get_available_bytes()
		mutex.unlock()
		if available_bytes:
			mutex.lock()
			var data_result = client.get_data(available_bytes)
			mutex.unlock()
			buff.append_array(data_result[1])
			var result = GameProtocol.decode(buff)
			var length = result[0]
			var proto_message = result[1]
			
			if proto_message == null:
				continue

			network_debug.log("Recived: %d %s %s" % [proto_message.type, proto_message.request_id, proto_message.message])

			mutex.lock()
			if pending_requests.has(proto_message.request_id):
				var pr = pending_requests[proto_message.request_id]
				pending_requests.erase(proto_message.request_id)
				mutex.unlock()
				print(proto_message.message)
				pr.response = proto_message.message
				pr.response_time = OS.get_ticks_msec()
				network_ping.call_deferred("set_text", "%d ms" % (pr.response_time - pr.request_time))
				if pr._timeout_timer != null:
					pr._timeout_timer.stop()
					remove_child(pr._timeout_timer)

				if pr._response_callback_obj != null and pr._response_callback_func != null:
					var response_args = [pr]
					if pr._response_args != null:
						for arg in pr._response_args:
							response_args.append(arg)
					pr._response_callback_obj.call_deferred(pr._response_callback_func, response_args)
			else:
				mutex.unlock()

			if len(buff) > length:
				buff = buff.subarray(length, len(buff) - 1)
			else:
				buff = PoolByteArray()

func set_auth_data(username):
	self.username = username

func auth():
	send({"name": username}, MessageTypes.AUTHENTICATE)

func start_thread(host, port):
	print("Spawned new network thread")
	thread = Thread.new()
	thread.start(self, "_start", [host, port])

func _start(data):
	mutex.lock()
	host = data[0]
	port = data[1]

	client = StreamPeerTCP.new()
	client.connect_to_host(host, port)

	if client.get_status() != 2:
		yield(get_tree().create_timer(1), "timeout")
		self.call_deferred("reconnect")
		return
	mutex.unlock()
	
	network_ok()
	auth()

	KeepAliveTimer = Timer.new()
	add_child(KeepAliveTimer)
	KeepAliveTimer.set_wait_time(KEEP_ALIVE_INVERVAL)
	KeepAliveTimer.connect("timeout", self, "_on_KeepAliveTimer_timeout")
	KeepAliveTimer.start()

	recv_message_loop()
	
func stop():
	_kill_thread = true
	thread.wait_to_finish()
	_kill_thread = false
	client.disconnect_from_host()
	network_error()

func reconnect():
	_kill_thread = true
	thread.wait_to_finish()
	_kill_thread = false
	print("Killed network thread")

	for pr in pending_requests.values():
		pr._timeout_timer.stop()
		remove_child(pr._timeout_timer)

	mutex = Mutex.new()
	pending_requests.clear()
	if KeepAliveTimer != null:
		KeepAliveTimer.stop()
		remove_child(KeepAliveTimer)
		KeepAliveTimer = null

	self.start_thread(host, port)

func _on_KeepAliveTimer_timeout():
	mutex.lock()
	KeepAlive.send()
	mutex.unlock()