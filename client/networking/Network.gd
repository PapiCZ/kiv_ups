extends Node

var client = null
var mutex = null
var thread = null
var pending_requests = {}

var PendingRequest = preload("res://networking/PendingRequest.gd")
var GameProtocol = preload("res://networking/GameProtocol.gd").new()
var ProtoMessage = preload("res://networking/ProtoMessage.gd")
var Utils = preload("res://networking/Utils.gd")
onready var status_sprite = get_tree().get_root().get_node("Game/NetworkStatus/Sprite")
onready var network_debug = get_tree().get_root().get_node("Game/NetworkDebug")

func _ready():
	mutex = Mutex.new()

func network_ok():
	status_sprite.call_deferred("set_texture", preload("res://sprites/network_ok.png"))
	
func network_error():
	status_sprite.call_deferred("set_texture", preload("res://sprites/network_error.png"))

func send(message, type, \
	response_callback_obj=null, response_callback_func=null, response_args=null, \
	timeout_callback_obj=null, timeout_callback_func=null, timeout_args=null):
	var proto_message = ProtoMessage.new()
	proto_message.message = message
	proto_message.type = type
	proto_message.request_id = Utils.random_request_id()

	var data = GameProtocol.encode(proto_message)
	network_debug.log("Sent %s" % data.get_string_from_ascii())

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
		pr._timeout_timer.set_wait_time(pr._timeout)
		var _timeout_args = [pr]
		if pr._timeout_args != null:
			for arg in pr._timeout_args:
				_timeout_args.append(arg)
		if pr._timeout_callback_obj != null and pr._timeout_callback_func != null:
			pr._timeout_timer.connect("timeout", pr._timeout_callback_obj, pr._timeout_callback_func, _timeout_args)
		pr._timeout_timer.connect("timeout", self, "_on_RequestTimeoutTimer_timeout", _timeout_args)
		pr._timeout_timer.start()

		add_child(pr._timeout_timer)

	mutex.lock()
	pending_requests[pr.id] = pr
	client.put_data(data)
	mutex.unlock()
	
	return data

func _on_RequestTimeoutTimer_timeout(pr):
	remove_child(pr._timeout_timer)

	mutex.lock()
	pending_requests.erase(pr.id)
	mutex.unlock()

func recv_message_loop():
	var buff = PoolByteArray()

	while true:
		var data_result = client.get_data(1) # Block until data are received
		buff.append_array(data_result[1])

		var available_bytes = client.get_available_bytes()
		if available_bytes:
			data_result = client.get_data(available_bytes)
			buff.append_array(data_result[1])
			var result = GameProtocol.decode(buff)
			var length = result[0]
			var proto_message = result[1]

			network_debug.log("Recived: %d %s %s" % [proto_message.type, proto_message.request_id, proto_message.message])

			var pr = null
			mutex.lock()
			if pending_requests.has(proto_message.request_id):
				pr = pending_requests[proto_message.request_id]
				pending_requests.erase(proto_message.request_id)
			mutex.unlock()

			if pr != null:
				pr.response = proto_message.message
				pr.response_time = OS.get_ticks_msec()
				if pr._timeout_timer != null:
					remove_child(pr._timeout_timer)

				if pr._response_callback_obj != null and pr._response_callback_func != null:
					var response_args = [pr]
					if pr._response_args != null:
						for arg in pr._response_args:
							response_args.append(arg)
					pr._response_callback_obj.call_deferred(pr._response_callback_func, response_args)


			if len(buff) > length:
				buff = buff.subarray(length, len(buff) - 1)
			else:
				buff = PoolByteArray()
	
func start_thread(host, port):
	thread = Thread.new()
	thread.start(self, "_start", [host, port])

func _start(data):
	var host = data[0]
	var port = data[1]

	client = StreamPeerTCP.new()
	client.connect_to_host(host, port)
	network_ok()
	network_debug.log("Connected to %s:%d" % [host, port])
	recv_message_loop()
	
func stop():
	client.disconnect_from_host()
	network_error()
