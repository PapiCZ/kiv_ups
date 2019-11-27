package actions

import (
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp"
	"kiv_ups_server/net/tcp/protocol"
)

type ActionResponse struct {
	ServerMessage tcp.ServerMessage
	Targets       []interfaces.Player
}

type Action interface {
	Process(interfaces.MasterServer, interfaces.PlayerMessage) ActionResponse
	GetPlayerContexts() []interfaces.PlayerContext
	GetMessage() protocol.Message
}

type ActionDefinition struct {
	actionMap map[interfaces.PlayerContext]map[protocol.MessageType]Action
}

func NewDefinition() ActionDefinition {
	return ActionDefinition{actionMap: make(map[interfaces.PlayerContext]map[protocol.MessageType]Action)}
}

func (ad *ActionDefinition) Register(action Action) {
	for _, ctx := range action.GetPlayerContexts() {
		playerCtx := ctx

		if _, ok := ad.actionMap[playerCtx]; !ok {
			ad.actionMap[playerCtx] = make(map[protocol.MessageType]Action)
		}

		ad.actionMap[playerCtx][action.GetMessage().GetTypeId()] = action
	}
}

func (ad *ActionDefinition) GetAction(messageType protocol.MessageType, context interfaces.PlayerContext) Action {
	return ad.actionMap[context][messageType]
}
