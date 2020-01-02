package actions

import (
	interfaces "kiv_ups_server/masterserver/interfaces"
	protocol "kiv_ups_server/net/tcp/protocol"
)

/*
###############################################################
# This code is generated by `go run ./scripts/generate.go`.   #
# If you want to generate the contents of this file, go to    #
# the server directory and run the command.                   #
#                                                             #
# DON'T TOUCH IT DIRECTLY! YOU WILL SUFFER!                   #
###############################################################
*/
const (
	DefaultContext      = interfaces.PlayerContext(0)
	LobbyContext        = interfaces.PlayerContext(1)
	InGameContext       = interfaces.PlayerContext(2)
	LoggedInMenuContext = interfaces.PlayerContext(3)
)

type KeepAliveAction struct{}

func (a KeepAliveAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{DefaultContext, LobbyContext, InGameContext, LoggedInMenuContext}
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

type JoinLobbyAction struct{}

func (a JoinLobbyAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{LoggedInMenuContext}
}

func (a JoinLobbyAction) GetMessage() protocol.Message {
	return protocol.JoinLobbyMessage{}
}

type ListLobbyPlayersAction struct{}

func (a ListLobbyPlayersAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{LobbyContext}
}

func (a ListLobbyPlayersAction) GetMessage() protocol.Message {
	return protocol.ListLobbyPlayersMessage{}
}

type StartGameAction struct{}

func (a StartGameAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{LobbyContext}
}

func (a StartGameAction) GetMessage() protocol.Message {
	return protocol.StartGameMessage{}
}

type GameReconnectAvailableAction struct{}

func (a GameReconnectAvailableAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{LoggedInMenuContext}
}

func (a GameReconnectAvailableAction) GetMessage() protocol.Message {
	return protocol.GameReconnectAvailableMessage{}
}

type ReconnectAction struct{}

func (a ReconnectAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{LoggedInMenuContext}
}

func (a ReconnectAction) GetMessage() protocol.Message {
	return protocol.ReconnectMessage{}
}

type LeaveGameAction struct{}

func (a LeaveGameAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{LoggedInMenuContext, InGameContext}
}

func (a LeaveGameAction) GetMessage() protocol.Message {
	return protocol.LeaveGameMessage{}
}

type LeaveLobbyAction struct{}

func (a LeaveLobbyAction) GetPlayerContexts() []interfaces.PlayerContext {
	return []interfaces.PlayerContext{LobbyContext}
}

func (a LeaveLobbyAction) GetMessage() protocol.Message {
	return protocol.LeaveLobbyMessage{}
}

func RegisterAllActions(actionDefinition *ActionDefinition) {
	actionDefinition.Register(KeepAliveAction{})
	actionDefinition.Register(AuthenticateAction{})
	actionDefinition.Register(CreateLobbyAction{})
	actionDefinition.Register(DeleteLobbyAction{})
	actionDefinition.Register(ListLobbiesAction{})
	actionDefinition.Register(JoinLobbyAction{})
	actionDefinition.Register(ListLobbyPlayersAction{})
	actionDefinition.Register(StartGameAction{})
	actionDefinition.Register(GameReconnectAvailableAction{})
	actionDefinition.Register(ReconnectAction{})
	actionDefinition.Register(LeaveGameAction{})
	actionDefinition.Register(LeaveLobbyAction{})
}
