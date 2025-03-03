extends Node

const DELIMITER_CHARACTER = 124 # "|"
var ProtoMessage = preload("res://networking/ProtoMessage.gd")

func encode(proto_message):
	var data = PoolByteArray()
	var json = JSON.print(proto_message.message).to_utf8()
	
	data.append(DELIMITER_CHARACTER)
	data.append_array(String(proto_message.type).to_utf8())
	data.append(DELIMITER_CHARACTER)
	data.append_array(String(len(json)).to_utf8())
	data.append(DELIMITER_CHARACTER)
	data.append_array(proto_message.request_id.to_utf8())
	data.append(DELIMITER_CHARACTER)
	data.append_array(json)
	
	return data

func decode(buff):
	var proto_message = null
	var offset = 0
	var buff_len = len(buff)

	# Read delimiter
	if offset >= buff_len:
		return [null, null]
	if buff[offset] == DELIMITER_CHARACTER:
		offset += 1
	
	# Read message type
	if offset >= buff_len:
		return [null, null]
	var ascii_number = read_ascii_number_until_delimiter(buff, offset)
	if ascii_number == null or ascii_number[0] == null:
		return [null, null]
	offset += ascii_number[0]
	var type = int(ascii_number[1].get_string_from_utf8())

	# Read delimiter
	if offset >= buff_len:
		return [null, null]
	if buff[offset] == DELIMITER_CHARACTER:
		offset += 1

	# Read JSON len
	if offset >= buff_len:
		return [null, null]
	var json_len_ascii = read_ascii_number_until_delimiter(buff, offset)
	if json_len_ascii == null or json_len_ascii[0] == null:
		return [null, null]
	offset += json_len_ascii[0]
	var json_len = int(json_len_ascii[1].get_string_from_utf8())

	# Read delimiter
	if offset >= buff_len:
		return [null, null]
	if buff[offset] == DELIMITER_CHARACTER:
		offset += 1

	# Read request ID
	if offset >= buff_len:
		return [null, null]
	var request_id_ascii = read_ascii_word_until_delimiter(buff, offset)
	if request_id_ascii == null or request_id_ascii[0] == null:
		return [null, null]
	offset += request_id_ascii[0]
	var request_id = request_id_ascii[1].get_string_from_utf8()

	# Read delimiter
	if offset >= buff_len:
		return [null, null]
	if buff[0] == DELIMITER_CHARACTER:
		offset += 1

	# Read JSON
	if offset >= buff_len:
		return [null, null]
	if offset + json_len - 1 >= len(buff):
		return [null, null]

	var json_result = JSON.parse(buff.subarray(offset, offset + json_len - 1).get_string_from_utf8())
	if json_result.error == OK:
		proto_message = ProtoMessage.new()
		proto_message.message = json_result.result
		proto_message.type = type
		proto_message.request_id = request_id
		return [offset + json_len, proto_message]
	else:
		return [null, null]

func read_ascii_number_until_delimiter(buff, start=0):
	var ascii_number_buff = PoolByteArray()

	var i = start
	while true:
		if i >= len(buff):
			return [null, null]

		if buff[i] >= 48 and buff[i] <= 57:
			# Between 0 and 9
			ascii_number_buff.append(buff[i])
		elif buff[i] == DELIMITER_CHARACTER:
			return [i - start, ascii_number_buff]
		else:
			return [null, null]

		i += 1

func read_ascii_word_until_delimiter(buff, start=0):
	var ascii_word_buff = PoolByteArray()

	var i = start
	while true:
		if i >= len(buff):
			return [null, null]

		if (buff[i] >= 48 and buff[i] <= 57) or (buff[i] >= 97 and buff[i] <= 122) or (buff[i] >= 65 and buff[i] <= 90):
			# 0-9 OR a-z OR A-Z
			ascii_word_buff.append(buff[i])
		elif buff[i] == DELIMITER_CHARACTER:
			return [i - start, ascii_word_buff]
		else:
			return [null, null]
		
		i += 1