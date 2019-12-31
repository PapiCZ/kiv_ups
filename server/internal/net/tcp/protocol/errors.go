package protocol

import (
	"fmt"
)

type UnexpectedCharacter struct {
	Character byte
}

func (uc UnexpectedCharacter) Error() string {
	return fmt.Sprintf("got an unexpected character %c", uc.Character)
}

type InvalidType struct {
	Type int
}

func (it InvalidType) Error() string {
	return fmt.Sprintf("got invalid message type %c", it.Type)
}
