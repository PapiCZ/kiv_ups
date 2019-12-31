package gameserver

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	tree2 "kiv_ups_server/internal/masterserver/gameserver/tree"
	nodes2 "kiv_ups_server/internal/masterserver/gameserver/tree/nodes"
	interfaces2 "kiv_ups_server/internal/masterserver/interfaces"
	"kiv_ups_server/internal/net/tcp"
	"kiv_ups_server/internal/net/tcp/protocol"
	"reflect"
	"time"
)

const MaxScore = 10000

// GameServer is a structure that handles all data about ingame instance.
type GameServer struct {
	Tps                 int
	tickDuration        int64
	startTime           time.Time
	Players             []interfaces2.Player
	DisconnectedPlayers []interfaces2.Player
	RequestMessageChan  chan interfaces2.PlayerMessage
	GameTree            tree2.Node
	Run                 bool
}

// NewGameServer creates and initializes new
func NewGameServer() GameServer {
	return GameServer{
		Tps:                30,
		Players:            make([]interfaces2.Player, 0),
		RequestMessageChan: make(chan interfaces2.PlayerMessage, 128),
		GameTree: tree2.Node{
			Parent:   nil,
			Children: nil,
			Value:    &nodes2.RootNode{},
		},
	}
}

// getElapsedMillis get number of milliseconds that has elapsed from GameServer
// start
func (gs *GameServer) getElapsedMillis() int64 {
	return time.Since(gs.startTime).Milliseconds()
}

// readRequestChanToSlice converts channel data to slice. The purpose is to
// make it possible to read every message multiple times.
func (gs *GameServer) readRequestChanToSlice() []interfaces2.PlayerMessage {
	requests := make([]interfaces2.PlayerMessage, 0, 10)

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

// buildTree adds all required nodes to game tree
func (gs *GameServer) buildTree() {
	for _, player := range gs.GetPlayers() {
		score := nodes2.Score{
			Player: player,
		}
		gs.GameTree.AddGameNodes(&nodes2.Spaceship{Player: player}, &score)
	}
	gs.GameTree.AddGameNodes(&nodes2.AsteroidWrapper{})
}

// initTree initializes all nodes in tree
func (gs *GameServer) initTree() {
	gs.GameTree.Init()

	for _, node := range gs.GameTree.GetAllChildren() {
		node.Init()
		node.Value.Init(node)
	}
}

// Start initializes game tree and other values, and starts gameloop
// that process all the magic in the game
func (gs *GameServer) Start() {
	log.Infoln("Starting game server...")

	gs.Run = true
	gs.buildTree()
	gs.initTree()

	gs.startTime = time.Now()
	gs.tickDuration = int64(1000 / gs.Tps)

	nextGameTickTime := gs.getElapsedMillis()

	playerMessagesBuffer := make([]interfaces2.PlayerMessage, 0)
	lastElapsedMillis := int64(0)

	// gameloop
	for gs.Run {
		playerMessagesBuffer = append(playerMessagesBuffer, gs.readRequestChanToSlice()...)

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
			playerMessagesBuffer = make([]interfaces2.PlayerMessage, 0)
		}
	}
}

// ManageGame is primarily used to check whether gameloop should stop
func (gs *GameServer) ManageGame() {
	// Check if someone has over MaxScore
	for _, node := range gs.GameTree.FindAllChildrenByType("score") {
		score := node.Value.(*nodes2.Score)
		if score.Score > MaxScore {
			// Build score summary
			scoreSummary := make([]PlayerScoreSummary, 0)
			for _, _node := range gs.GameTree.FindAllChildrenByType("score") {
				_score := _node.Value.(*nodes2.Score)
				scoreSummary = append(scoreSummary, PlayerScoreSummary{
					PlayerName: _score.PlayerName,
					Score:      _score.Score,
					Winner:     score.Player == _score.Player,
				})
			}

			// Reset certain variables and send score summary to all players
			for _, player := range gs.Players {
				player.SetGameServer(nil)
				player.SetLoggedInMenuContext()
				player.SetConnectedLobby(nil)
				_ = player.GetTCPClient().Send(protocol.ProtoMessage{
					Message: &tcp.ServerMessage{
						Data:    protocol.GameEndMessage{ScoreSummary: scoreSummary},
						Status:  true,
						Message: "",
					},
					RequestId: "",
				})
			}

			// Shutdown gameloop
			gs.Run = false
			return
		}
	}

	// Check if is any player active
	active := false
	for _, player := range gs.Players {
		if player.IsConnected() {
			active = true
			break
		}
	}

	if !active {
		gs.Run = false
	}
}

// shareState sends current state of game tree to all players
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

// AddPlayer adds player to game server
func (gs *GameServer) AddPlayer(player interfaces2.Player) {
	gs.Players = append(gs.Players, player)
	player.SetLoggedInMenuContext()
	player.SetGameServer(gs)
}

// RemovePlayer removes player from game server
func (gs *GameServer) RemovePlayer(player interfaces2.Player) {
	for i, _player := range gs.Players {
		if player == _player {
			gs.Players = append(gs.Players[:i], gs.Players[i+1:]...)
		}
	}

	// remove spaceship if game has already started
	if gs.Run {
		for _, node := range gs.GameTree.FindAllChildrenByType("spaceship") {
			spaceship := node.Value.(*nodes2.Spaceship)

			if spaceship.Player == player {
				node.Destroy()
				break
			}
		}
	}

	// remove score if game has already started
	if gs.Run {
		for _, node := range gs.GameTree.FindAllChildrenByType("score") {
			score := node.Value.(*nodes2.Score)

			fmt.Printf("%#v\n", score)

			if score.Player == player {
				node.Destroy()
				break
			}
		}
	}
}

// GetPlayers is getter for players
func (gs *GameServer) GetPlayers() []interfaces2.Player {
	return gs.Players
}

// GetRequestMessageChan is getter for request message channel
func (gs *GameServer) GetRequestMessageChan() chan interfaces2.PlayerMessage {
	return gs.RequestMessageChan
}

// AddDisconnectedPlayer adds player to slice of disconnected players
func (gs *GameServer) AddDisconnectedPlayer(player interfaces2.Player) {
	gs.DisconnectedPlayers = append(gs.DisconnectedPlayers, player)
}

// RemoveDisconnectedPlayer removes player from slice of disconnected players
func (gs *GameServer) RemoveDisconnectedPlayer(player interfaces2.Player) {
	for i, disconnectedPlayer := range gs.DisconnectedPlayers {
		if player == disconnectedPlayer {
			gs.DisconnectedPlayers = append(gs.DisconnectedPlayers[:i], gs.DisconnectedPlayers[i+1:]...)
			return
		}
	}
}

// IsPlayerDisconnected checks if player is in slice of disconnected players
func (gs *GameServer) IsPlayerDisconnected(player interfaces2.Player) bool {
	for _, disconnectedPlayer := range gs.DisconnectedPlayers {
		if player == disconnectedPlayer {
			return true
		}
	}

	return false
}

// Checks if gameloop is running
func (gs *GameServer) IsRunning() bool {
	return gs.Run
}

// FilterMessagesByTypes filters given messages by given types and
// returns slice of filtered messages
func FilterMessagesByTypes(playerMessages []interfaces2.PlayerMessage, types []protocol.Message) []interfaces2.PlayerMessage {
	filteredMessages := make([]interfaces2.PlayerMessage, 0)

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
	Winner     bool   `json:"winner"`
}
