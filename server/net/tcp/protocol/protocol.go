package protocol

import (
	"bufio"
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
)

const DelimiterCharacter = byte('|')

type MessageType int
type Message interface {
	GetTypeId() MessageType
}

type RequestId string

type ProtoMessage struct {
	Message
	RequestId
}

type Definition struct {
	types map[MessageType]Message
}

func NewDefinition() Definition {
	return Definition{types: make(map[MessageType]Message)}
}

func (md *Definition) Register(type_ Message) {
	md.types[type_.GetTypeId()] = type_
}

type Protocol interface {
	Encode(Message, io.Writer)
	Decode(io.Reader) Message
}

type GameProtocol struct {
	Def Definition
}

func (gp *GameProtocol) Encode(msg ProtoMessage, writer io.WriteCloser) {
	var err error
	buff := bufio.NewWriter(writer)

	// Write delimiter character
	_, err = buff.Write([]byte{DelimiterCharacter})
	if err != nil {
		log.Errorln(err)
	}

	// Write message type
	msgType := strconv.Itoa(int(msg.GetTypeId()))
	_, err = buff.Write([]byte(msgType))

	// Write delimiter character
	_, err = buff.Write([]byte{DelimiterCharacter})
	if err != nil {
		log.Errorln(err)
	}

	// Prepare JSON
	jsonBuff := &bytes.Buffer{}
	encoder := json.NewEncoder(jsonBuff)
	err = encoder.Encode(msg.Message)
	if err != nil {
		log.Errorln(err)
	}

	// Write JSON length
	jsonAsciiLen := strconv.Itoa(jsonBuff.Len() - 1)
	_, err = buff.Write([]byte(jsonAsciiLen))
	if err != nil {
		log.Errorln(err)
	}

	// Write delimiter character
	_, err = buff.Write([]byte{DelimiterCharacter})
	if err != nil {
		log.Errorln(err)
	}

	// Write request ID
	_, err = buff.Write([]byte(msg.RequestId))
	if err != nil {
		log.Errorln(err)
	}

	// Write delimiter character
	_, err = buff.Write([]byte{DelimiterCharacter})
	if err != nil {
		log.Errorln(err)
	}

	// Write JSON
	_, err = buff.Write(jsonBuff.Bytes()[:jsonBuff.Len() - 1])
	if err != nil {
		log.Errorln(err)
	}

	err = buff.Flush()
	if err != nil {
		log.Errorln(err)
	}
	_ = writer.Close()
}

func (gp *GameProtocol) InfiniteDecode(reader io.ReadCloser, msgChan chan *ProtoMessage) {
	// TODO: This goroutine will never die
	for {
		msg, err := gp.Decode(reader)
		if err != nil {
			if err == io.ErrClosedPipe {
				log.Traceln("Closing InfiniteDecode...")
				return
			} else {
				log.Errorln(err)
			}

			continue
		}

		msgChan <- msg
	}
}

func (gp *GameProtocol) Decode(reader io.Reader) (*ProtoMessage, error) {
	// TODO: Make it more strict
	// This shouldn't work.. |100|15|{"ping":"pong"} (missing request ID)
	var err error

	buffReader := bufio.NewReader(reader)

	// Just for sure
	err = gp.flushInvalidBytes(buffReader)
	if err != nil {
		if err == io.ErrClosedPipe {
			return nil, err
		} else {
			log.Errorln(err)
		}
	}

	// Flush "|" character
	err = gp.flushUntilDelimiter(buffReader)
	if err != nil {
		if err == io.ErrClosedPipe {
			return nil, err
		} else {
			log.Errorln(err)
		}
	}

	// Read message type and flush delimiter
	typeInt, err := gp.readAsciiNumberUntilDelimiter(buffReader)
	if err != nil {
		if err == io.ErrClosedPipe {
			return nil, err
		} else {
			log.Errorln(err)
		}
	}

	// Read JSON length and flush delimiter
	jsonLenInt, err := gp.readAsciiNumberUntilDelimiter(buffReader)
	if err != nil {
		if err == io.ErrClosedPipe {
			return nil, err
		} else {
			log.Errorln(err)
		}
	}

	// Read request ID
	requestId, err := gp.readAsciiWordUntilDelimiter(buffReader)
	if err != nil {
		if err == io.ErrClosedPipe {
			return nil, err
		} else {
			log.Warnln(err)
		}
	}

	// Read JSON
	limitedReader := &io.LimitedReader{R: buffReader, N: int64(jsonLenInt)}
	decoder := json.NewDecoder(limitedReader)

	// Get message by type
	target, ok := gp.Def.types[MessageType(typeInt)]

	var msg interface{}
	if !ok {
		return nil, InvalidType{Type: typeInt}
	}

	msg = reflect.New(reflect.TypeOf(target)).Interface()
	err = decoder.Decode(&msg)
	if err != nil {
		return nil, err
	}

	// Flush rest of the buffer
	_, err = ioutil.ReadAll(limitedReader)
	if err != nil {
		log.Errorln(err)
	}

	return &ProtoMessage{
		Message:     msg.(Message),
		RequestId:   RequestId(requestId),
	}, nil
}

func (gp *GameProtocol) flushInvalidBytes(reader *bufio.Reader) error {
	for {
		delimiterBuff := make([]byte, 1)
		_, err := io.ReadFull(reader, delimiterBuff)

		if err != nil {
			return err
		}
		if delimiterBuff[0] == DelimiterCharacter {
			err := reader.UnreadByte()
			if err != nil {
				return err
			}
			return nil
		}
	}
}

func (gp *GameProtocol) flushUntilDelimiter(reader io.Reader) error {
	for {
		delimiterBuff := make([]byte, 1)
		_, err := io.ReadFull(reader, delimiterBuff)
		if err != nil {
			return err
		}
		if delimiterBuff[0] != DelimiterCharacter {
			return UnexpectedCharacter{Character: delimiterBuff[0]}
		}

		return nil
	}
}

func (gp *GameProtocol) flushDelimiter(reader io.Reader) error {
	delimiterBuff := make([]byte, 1)
	_, err := io.ReadFull(reader, delimiterBuff)
	if err != nil {
		return err
	}
	if delimiterBuff[0] != DelimiterCharacter {
		return UnexpectedCharacter{Character: delimiterBuff[0]}
	}

	return nil
}

func (gp *GameProtocol) readAsciiNumberUntilDelimiter(reader *bufio.Reader) (int, error) {
	var asciiNumberBuff []byte
	asciiDigitBuff := make([]byte, 1)

	for {
		_, err := reader.Read(asciiDigitBuff)
		if err != nil {
			return 0, err
		}

		if asciiDigitBuff[0] >= '0' && asciiDigitBuff[0] <= '9' {
			asciiNumberBuff = append(asciiNumberBuff, asciiDigitBuff[0])
		} else if asciiDigitBuff[0] == DelimiterCharacter {
			return strconv.Atoi(string(asciiNumberBuff))
		} else {
			_ = reader.UnreadByte()
			return 0, UnexpectedCharacter{Character: asciiDigitBuff[0]}
		}
	}
}

func (gp *GameProtocol) readAsciiWordUntilDelimiter(reader *bufio.Reader) (string, error) {
	var asciiWordBuff []byte
	asciiDigitBuff := make([]byte, 1)

	for {
		_, err := reader.Read(asciiDigitBuff)
		if err != nil {
			return "", err
		}

		if (asciiDigitBuff[0] >= 'a' && asciiDigitBuff[0] <= 'z') ||
			(asciiDigitBuff[0] >= 'A' && asciiDigitBuff[0] <= 'Z') ||
			(asciiDigitBuff[0] >= '0' && asciiDigitBuff[0] <= '9') {
			asciiWordBuff = append(asciiWordBuff, asciiDigitBuff[0])
		} else if asciiDigitBuff[0] == DelimiterCharacter {
			return string(asciiWordBuff), nil
		} else {
			_ = reader.UnreadByte()
			return "", UnexpectedCharacter{Character: asciiDigitBuff[0]}
		}
	}
}
