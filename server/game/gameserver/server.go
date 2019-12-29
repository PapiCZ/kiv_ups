package gameserver

import (
	log "github.com/sirupsen/logrus"
	"kiv_ups_server/game/gameserver/tree"
	"kiv_ups_server/game/gameserver/tree/nodes"
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp"
	"kiv_ups_server/net/tcp/protocol"
	"reflect"
	"time"
)

type GameServer struct {
	Tps                 int
	tickDuration        int64
	startTime           time.Time
	Players             []interfaces.Player
	DisconnectedPlayers []interfaces.Player
	RequestMessageChan  chan interfaces.PlayerMessage
	GameTree            tree.Node
	Run                 bool
}

func NewGameServer() GameServer {
	return GameServer{
		Tps:                30,
		Players:            make([]interfaces.Player, 0),
		RequestMessageChan: make(chan interfaces.PlayerMessage, 128),
		GameTree: tree.Node{
			Parent:   nil,
			Children: nil,
			Value:    &nodes.RootNode{},
		},
	}
}

func (gs *GameServer) getElapsedMillis() int64 {
	return time.Since(gs.startTime).Milliseconds()
}

func (gs *GameServer) readRequestChanToList() []interfaces.PlayerMessage {
	requests := make([]interfaces.PlayerMessage, 0, 10)

	done := false
	for {
		select {
		case x, ok := <-gs.RequestMessageChan:
			if ok {
				requests = append(requests, x)
			} else {
				done = true
			}
		default:
			done = true
		}

		if done {
			break
		}
	}

	return requests
}

func (gs *GameServer) buildTree() {
	for _, player := range gs.GetPlayers() {
		score := nodes.Score{
			Player: player,
		}
		gs.GameTree.AddGameNodes(&nodes.Spaceship{Player: player}, &score)
	}
	gs.GameTree.AddGameNodes(&nodes.AsteroidWrapper{})
}

func (gs *GameServer) initTree() {
	gs.GameTree.Init()

	for _, node := range gs.GameTree.GetAllChildren() {
		node.Init()
		node.Value.Init(node)
	}
}

func (gs *GameServer) Start() {
	log.Infoln("Starting game server...")

	gs.Run = true
	gs.buildTree()
	gs.initTree()

	gs.startTime = time.Now()
	gs.tickDuration = int64(1000 / gs.Tps)

	nextGameTickTime := gs.getElapsedMillis()

	playerMessagesBuffer := make([]interfaces.PlayerMessage, 0)
	lastElapsedMillis := int64(0)

	// gameloop
	for gs.Run {
		playerMessagesBuffer = append(playerMessagesBuffer, gs.readRequestChanToList()...)

		for gs.getElapsedMillis() > nextGameTickTime {
			// update game
			elapsedMillis := gs.getElapsedMillis()
			for _, node := range gs.GameTree.GetAllChildren() {
				gameNode := node.Value
				gameNode.Process(
					gameNode.Filter(FilterMessagesByTypes(playerMessagesBuffer, gameNode.ListenMessages())),
					float64(elapsedMillis-lastElapsedMillis)/1000.0,
				)
			}

			for _, player := range gs.Players {
				if !player.IsConnected() && !gs.IsPlayerDisconnected(player) {
					gs.AddDisconnectedPlayer(player)

					// Send notification to all connected players
					for _, _player := range gs.Players {
						if _player.IsConnected() {
							_ = _player.GetTCPClient().Send(protocol.ProtoMessage{
								Message: &tcp.ServerMessage{
									Data:    protocol.PlayerDisconnectedMessage{PlayerName: player.GetName()},
									Status:  true,
									Message: "",
								},
								RequestId: "",
							})
						}
					}
				}
			}

			gs.ManageGame()
			gs.shareState()
			lastElapsedMillis = elapsedMillis
			nextGameTickTime += gs.tickDuration
			playerMessagesBuffer = make([]interfaces.PlayerMessage, 0)
		}
	}
}

func (gs *GameServer) ManageGame() {
	// Check if someone has over 50000
	for _, node := range gs.GameTree.FindAllChildrenByType("score") {
		score := node.Value.(*nodes.Score)
		if score.Score > 50000 {
			scoreSummary := make([]PlayerScoreSummary, 0)
			for _, _node := range gs.GameTree.FindAllChildrenByType("score") {
				_score := _node.Value.(*nodes.Score)
				scoreSummary = append(scoreSummary, PlayerScoreSummary{
					PlayerName: _score.PlayerName,
					Score:      _score.Score,
					Winner:     score.Player == _score.Player,
				})
			}

			for _, player := range gs.Players {
				player.SetGameServer(nil)
				_ = player.GetTCPClient().Send(protocol.ProtoMessage{
					Message: &tcp.ServerMessage{
						Data:    protocol.GameEndMessage{ScoreSummary: scoreSummary},
						Status:  true,
						Message: "",
					},
					RequestId: "",
				})
			}

			gs.Run = false
			return
		}
	}
}

func (gs *GameServer) shareState() {
	for _, player := range gs.GetPlayers() {
		err := player.GetTCPClient().Send(protocol.ProtoMessage{
			Message: &tcp.ServerMessage{
				Status:  true,
				Message: "",
				Data:    &protocol.UpdateStateMessage{GameTree: gs.GameTree},
			},
			RequestId: protocol.RequestId(""),
		})

		if err != nil {
			log.Errorln("Share state error:", err)
		}
	}
}

func (gs *GameServer) AddPlayer(player interfaces.Player) {
	gs.Players = append(gs.Players, player)
	player.SetGameServer(gs)
}

func (gs *GameServer) GetPlayers() []interfaces.Player {
	return gs.Players
}

func (gs *GameServer) GetRequestMessageChan() chan interfaces.PlayerMessage {
	return gs.RequestMessageChan
}

func (gs *GameServer) AddDisconnectedPlayer(player interfaces.Player) {
	gs.DisconnectedPlayers = append(gs.DisconnectedPlayers, player)
}

func (gs *GameServer) RemoveDisconnectedPlayer(player interfaces.Player) {
	for i, disconnectedPlayer := range gs.DisconnectedPlayers {
		if player == disconnectedPlayer {
			gs.DisconnectedPlayers = append(gs.DisconnectedPlayers[:i], gs.DisconnectedPlayers[i+1:]...)
			return
		}
	}
}

func (gs *GameServer) IsPlayerDisconnected(player interfaces.Player) bool {
	for _, disconnectedPlayer := range gs.DisconnectedPlayers {
		if player == disconnectedPlayer {
			return true
		}
	}

	return false
}

func FilterMessagesByTypes(playerMessages []interfaces.PlayerMessage, types []protocol.Message) []interfaces.PlayerMessage {
	filteredMessages := make([]interfaces.PlayerMessage, 0)

	for _, playerMessage := range playerMessages {
		for _, _type := range types {
			if reflect.ValueOf(playerMessage.GetMessage().Message).Type() == reflect.ValueOf(_type).Type() {
				filteredMessages = append(filteredMessages, playerMessage)
			}
		}
	}

	return playerMessages
}

type PlayerScoreSummary struct {
	PlayerName string `json:"player_name"`
	Score      int    `json:"score"`
	Winner     bool
}
