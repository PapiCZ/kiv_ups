extends VBoxContainer

var score_summary = {}

func _load():
	for ss in score_summary:
		if ss.winner:
			$VBoxContainer/PlayerNameLabel.text = ss.player_name

func _on_MainMenu_pressed():
	Menu.reset_all()
	Menu.go(Menu.MENU_LEVEL.MAIN)
