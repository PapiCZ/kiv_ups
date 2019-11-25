extends Node

static func random_request_id(length=5):
	var charset = [
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 
		'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 
		'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'
	]

	var random_str = ""
	randomize()

	#warning-ignore:unused_variable
	for i in range(length):
		random_str += charset[randi() % (len(charset))]

	return random_str

