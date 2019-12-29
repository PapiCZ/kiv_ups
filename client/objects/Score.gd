extends Node2D

var score = 0
var player_name = ""

func update_player_name(player_name):
	self.player_name = player_name
	$ProgressBar/PlayerNameLabel.text = self.player_name

func update_score(score):
	self.score = score
	$ProgressBar.value = self.score
	$ProgressBar/ScoreLabel.text = str(self.score)

func set_disconnected():
	$ProgressBar/PlayerNameLabel.text = "X " + $ProgressBar/PlayerNameLabel.text
	$ProgressBar/PlayerNameLabel.add_color_override("font_color", Color(1, 0, 0))

func set_connected():
	$ProgressBar/PlayerNameLabel.text.erase(0, 2)
	$ProgressBar/PlayerNameLabel.set("custom_colors/font_color", null)
