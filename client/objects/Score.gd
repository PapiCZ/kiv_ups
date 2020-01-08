extends Node2D

var score = 0
var player_name = ""
var connected = true

func update_player_name(player_name):
	self.player_name = player_name
	$ProgressBar/PlayerNameLabel.text = self.player_name

func update_score(score):
	self.score = score
	$ProgressBar.value = self.score
	$ProgressBar/ScoreLabel.text = str(self.score)

func set_disconnected():
	if connected:
		connected = false
		print("CONNECTED!")
		$ProgressBar/PlayerNameLabel.text = "X " + self.player_name
		$ProgressBar/PlayerNameLabel.add_color_override("font_color", Color(1, 0, 0))

func set_connected():
	if not connected:
		print("DISCONNECTED!")
		connected = true
		$ProgressBar/PlayerNameLabel.text = self.player_name
		$ProgressBar/PlayerNameLabel.set("custom_colors/font_color", null)
