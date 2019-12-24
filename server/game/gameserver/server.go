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

type GamePlayer struct {
	Player interfaces.Player
}

type GameServer struct {
	Tps                int
	tickDuration       int64
	startTime          time.Time
	Players            []*GamePlayer
	RequestMessageChan chan interfaces.PlayerMessage
	GameTree           tree.Node
}

func NewGameServer() GameServer {
	return GameServer{
		Tps:                30,
		Players:            make([]*GamePlayer, 0),
		RequestMessageChan: make(chan interfaces.PlayerMessage, 128),
		GameTree: tree.Node{
			Parent:   nil,
			Children: nil,
			Value:    nodes.RootNode{},
		},
	}
}

func (gs *GameServer) GetElapsedMillis() int64 {
	return time.Since(gs.startTime).Milliseconds()
}

func (gs *GameServer) ReadRequestChanToList() []interfaces.PlayerMessage {
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
		gs.GameTree.AddGameNodes(&nodes.Spaceship{Player: player})
	}
}

func (gs *GameServer) Start() {
	log.Infoln("Starting game server...")

	gs.buildTree()

	gs.startTime = time.Now()
	gs.tickDuration = int64(1000 / gs.Tps)

	nextGameTickTime := gs.GetElapsedMillis()

	playerMessagesBuffer := make([]interfaces.PlayerMessage, 0)

	// gameloop
	for {
		playerMessagesBuffer = append(playerMessagesBuffer, gs.ReadRequestChanToList()...)

		for gs.GetElapsedMillis() > nextGameTickTime {
			// update game
			for _, node := range gs.GameTree.GetAllChildren() {
				gameNode := node.Value
				gameNode.Process(
					gameNode.Filter(FilterMessagesByTypes(playerMessagesBuffer, gameNode.ListenMessages())),
					time.Duration(gs.GetElapsedMillis()-nextGameTickTime)*time.Second,
				)
			}

			gs.ShareState()
			nextGameTickTime += gs.tickDuration
			playerMessagesBuffer = make([]interfaces.PlayerMessage, 0)
		}
	}
}

func (gs *GameServer) ShareState() {
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
	gs.Players = append(gs.Players, &GamePlayer{
		Player: player,
	})
	player.SetGameServer(gs)
}

func (gs *GameServer) GetPlayers() []interfaces.Player {
	players := make([]interfaces.Player, 0)

	for _, gamePlayer := range gs.Players {
		players = append(players, gamePlayer.Player)
	}

	return players
}

func (gs *GameServer) GetRequestMessageChan() chan interfaces.PlayerMessage {
	return gs.RequestMessageChan
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
