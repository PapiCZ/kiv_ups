extends Node

#warning-ignore-all:unused_class_variable

var id = 0
var request = null
var request_time = null
var response = null
var response_time = null

var _response_callback_obj = null
var _response_callback_func = null
var _response_args

var _timeout_callback_obj = null
var _timeout_callback_func = null
var _timeout_args
var _timeout = 1 # seconds
var _timeout_timer = null