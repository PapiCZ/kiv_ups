package actions

import (
	interfaces "kiv_ups_server/game/interfaces"
	protocol "kiv_ups_server/net/tcp/protocol"
)

// ##############################################################
// This code is generated by `go run ./generate`. DON'T TOUCH IT!
// ##############################################################

const (
	DefaultContext      = interfaces.PlayerContext(0)
	LoggedInMenuContext = interfaces.PlayerContext(1)
)

type KeepAliveAction struct{}

func (a KeepAliveAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{DefaultContext, LoggedInMenuContext}
}

func (a KeepAliveAction) GetMessage() protocol.Message {
	return protocol.KeepAliveMessage{}
}

type AuthenticateAction struct{}

func (a AuthenticateAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{DefaultContext}
}

func (a AuthenticateAction) GetMessage() protocol.Message {
	return protocol.AuthenticateMessage{}
}

type CreateLobbyAction struct{}

func (a CreateLobbyAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{LoggedInMenuContext}
}

func (a CreateLobbyAction) GetMessage() protocol.Message {
	return protocol.CreateLobbyMessage{}
}

type DeleteLobbyAction struct{}

func (a DeleteLobbyAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{LoggedInMenuContext}
}

func (a DeleteLobbyAction) GetMessage() protocol.Message {
	return protocol.DeleteLobbyMessage{}
}

type ListLobbiesAction struct{}

func (a ListLobbiesAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{LoggedInMenuContext}
}

func (a ListLobbiesAction) GetMessage() protocol.Message {
	return protocol.ListLobbiesMessage{}
}

func RegisterAllActions(actionDefinition *ActionDefinition) {
	actionDefinition.Register(KeepAliveAction{})
	actionDefinition.Register(AuthenticateAction{})
	actionDefinition.Register(CreateLobbyAction{})
	actionDefinition.Register(DeleteLobbyAction{})
	actionDefinition.Register(ListLobbiesAction{})
}
