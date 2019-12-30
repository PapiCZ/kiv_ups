package protocol

// ##############################################################
// This code is generated by `go run ./generate`. DON'T TOUCH IT!
// ##############################################################

type KeepAliveMessage struct {
	Ping string `json:"ping"`
}

func (m KeepAliveMessage) GetTypeId() MessageType {
	return 100
}

type ActionErrorMessage struct{}

func (m ActionErrorMessage) GetTypeId() MessageType {
	return 101
}

type AuthenticateMessage struct {
	Name string `json:"name"`
}

func (m AuthenticateMessage) GetTypeId() MessageType {
	return 200
}

type CreateLobbyMessage struct {
	Name         string `json:"name"`
	PlayersLimit int    `json:"players_limit"`
}

func (m CreateLobbyMessage) GetTypeId() MessageType {
	return 201
}

type CreatedLobbyResponseMessage struct{}

func (m CreatedLobbyResponseMessage) GetTypeId() MessageType {
	return 301
}

type DeleteLobbyMessage struct {
	Name string `json:"name"`
}

func (m DeleteLobbyMessage) GetTypeId() MessageType {
	return 202
}

type DeleteLobbyResponseMessage struct{}

func (m DeleteLobbyResponseMessage) GetTypeId() MessageType {
	return 302
}

type ListLobbiesMessage struct{}

func (m ListLobbiesMessage) GetTypeId() MessageType {
	return 203
}

type ListLobbiesResponseMessage struct {
	Lobbies interface{} `json:"lobbies"`
}

func (m ListLobbiesResponseMessage) GetTypeId() MessageType {
	return 303
}

type JoinLobbyMessage struct {
	Name string `json:"name"`
}

func (m JoinLobbyMessage) GetTypeId() MessageType {
	return 204
}

type JoinLobbyResponseMessage struct{}

func (m JoinLobbyResponseMessage) GetTypeId() MessageType {
	return 304
}

type PlayerLobbyJoinedMessage struct {
	PlayerName string `json:"player_name"`
}

func (m PlayerLobbyJoinedMessage) GetTypeId() MessageType {
	return 305
}

type ListLobbyPlayersMessage struct{}

func (m ListLobbyPlayersMessage) GetTypeId() MessageType {
	return 206
}

type ListLobbyPlayersResponseMessage struct {
	Players []string `json:"players"`
}

func (m ListLobbyPlayersResponseMessage) GetTypeId() MessageType {
	return 306
}

type StartGameMessage struct{}

func (m StartGameMessage) GetTypeId() MessageType {
	return 207
}

type StartGameResponseMessage struct{}

func (m StartGameResponseMessage) GetTypeId() MessageType {
	return 307
}

type LobbyPlayerConnectedMessage struct {
	Name string `json:"name"`
}

func (m LobbyPlayerConnectedMessage) GetTypeId() MessageType {
	return 308
}

type LobbyPlayerDisconnectedMessage struct {
	Name string `json:"name"`
}

func (m LobbyPlayerDisconnectedMessage) GetTypeId() MessageType {
	return 309
}

type GameEndMessage struct {
	ScoreSummary interface{} `json:"score_summary"`
}

func (m GameEndMessage) GetTypeId() MessageType {
	return 310
}

type GameReconnectAvailableMessage struct{}

func (m GameReconnectAvailableMessage) GetTypeId() MessageType {
	return 211
}

type GameReconnectAvailableResponseMessage struct {
	Available bool `json:"available"`
}

func (m GameReconnectAvailableResponseMessage) GetTypeId() MessageType {
	return 311
}

type ReconnectMessage struct{}

func (m ReconnectMessage) GetTypeId() MessageType {
	return 212
}

type ReconnectResponseMessage struct{}

func (m ReconnectResponseMessage) GetTypeId() MessageType {
	return 312
}

type LeaveGameMessage struct{}

func (m LeaveGameMessage) GetTypeId() MessageType {
	return 213
}

type LeaveGameResponseMessage struct{}

func (m LeaveGameResponseMessage) GetTypeId() MessageType {
	return 313
}

type LeaveLobbyMessage struct{}

func (m LeaveLobbyMessage) GetTypeId() MessageType {
	return 214
}

type LeaveLobbyResponseMessage struct{}

func (m LeaveLobbyResponseMessage) GetTypeId() MessageType {
	return 314
}

type PlayerMoveMessage struct {
	PlayerName string  `json:"player_name"`
	PosX       float64 `json:"pos_x"`
	PosY       float64 `json:"pos_y"`
	VelocityX  float64 `json:"velocity_x"`
	VelocityY  float64 `json:"velocity_y"`
	Rotation   float64 `json:"rotation"`
}

func (m PlayerMoveMessage) GetTypeId() MessageType {
	return 400
}

type ShootProjectileMessage struct{}

func (m ShootProjectileMessage) GetTypeId() MessageType {
	return 401
}

type UpdateStateMessage struct {
	GameTree interface{} `json:"game_tree"`
}

func (m UpdateStateMessage) GetTypeId() MessageType {
	return 500
}

type PlayerDisconnectedMessage struct {
	PlayerName string `json:"player_name"`
}

func (m PlayerDisconnectedMessage) GetTypeId() MessageType {
	return 501
}

func RegisterAllMessages(definition *Definition) {
	definition.Register(KeepAliveMessage{})
	definition.Register(ActionErrorMessage{})
	definition.Register(AuthenticateMessage{})
	definition.Register(CreateLobbyMessage{})
	definition.Register(CreatedLobbyResponseMessage{})
	definition.Register(DeleteLobbyMessage{})
	definition.Register(DeleteLobbyResponseMessage{})
	definition.Register(ListLobbiesMessage{})
	definition.Register(ListLobbiesResponseMessage{})
	definition.Register(JoinLobbyMessage{})
	definition.Register(JoinLobbyResponseMessage{})
	definition.Register(PlayerLobbyJoinedMessage{})
	definition.Register(ListLobbyPlayersMessage{})
	definition.Register(ListLobbyPlayersResponseMessage{})
	definition.Register(StartGameMessage{})
	definition.Register(StartGameResponseMessage{})
	definition.Register(LobbyPlayerConnectedMessage{})
	definition.Register(LobbyPlayerDisconnectedMessage{})
	definition.Register(GameEndMessage{})
	definition.Register(GameReconnectAvailableMessage{})
	definition.Register(GameReconnectAvailableResponseMessage{})
	definition.Register(ReconnectMessage{})
	definition.Register(ReconnectResponseMessage{})
	definition.Register(LeaveGameMessage{})
	definition.Register(LeaveGameResponseMessage{})
	definition.Register(LeaveLobbyMessage{})
	definition.Register(LeaveLobbyResponseMessage{})
	definition.Register(PlayerMoveMessage{})
	definition.Register(ShootProjectileMessage{})
	definition.Register(UpdateStateMessage{})
	definition.Register(PlayerDisconnectedMessage{})
}
