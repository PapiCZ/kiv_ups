package actions

import (
	interfaces2 "kiv_ups_server/internal/masterserver/interfaces"
	"kiv_ups_server/internal/net/tcp"
	"kiv_ups_server/internal/net/tcp/protocol"
)

type ActionResponse struct {
	ServerMessage tcp.ServerMessage
	Targets       []interfaces2.Player
}

type Action interface {
	Process(interfaces2.MasterServer, interfaces2.PlayerMessage) ActionResponse
	GetPlayerContexts() []interfaces2.PlayerContext
	GetMessage() protocol.Message
}

type ActionDefinition struct {
	actionMap map[interfaces2.PlayerContext]map[protocol.MessageType]Action
}

func NewDefinition() ActionDefinition {
	return ActionDefinition{actionMap: make(map[interfaces2.PlayerContext]map[protocol.MessageType]Action)}
}

// Register registers given action according to message type ID and message contexts
func (ad *ActionDefinition) Register(action Action) {
	for _, ctx := range action.GetPlayerContexts() {
		playerCtx := ctx

		if _, ok := ad.actionMap[playerCtx]; !ok {
			ad.actionMap[playerCtx] = make(map[protocol.MessageType]Action)
		}

		ad.actionMap[playerCtx][action.GetMessage().GetTypeId()] = action
	}
}

// GetAction returns action according to given message type and context
func (ad *ActionDefinition) GetAction(messageType protocol.MessageType, context interfaces2.PlayerContext) Action {
	return ad.actionMap[context][messageType]
}
