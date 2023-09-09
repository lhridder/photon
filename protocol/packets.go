package protocol

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"strings"
)

type Packet struct {
	ID   byte
	Data *[]byte
}

// Handshake

const (
	HandshakeID             = 0x00
	HandshakeForgeSeparator = "\x00"
)

type Handshake struct {
	Protoversion int32
	ServerAddr   string
	ServerPort   uint16
	NextState    int32
}

func (p *Packet) ReadHandshake() (*Handshake, error) {
	if p.ID != HandshakeID {
		return nil, ErrPacketID
	}

	protoVersion, err := p.readVarint()
	if err != nil {
		return nil, ErrPacketParse
	}

	serverAddr, err := p.readString()
	if err != nil {
		return nil, ErrPacketParse
	}

	serverAddr = strings.Split(serverAddr, HandshakeForgeSeparator)[0]
	serverAddr = strings.Trim(serverAddr, ".")
	serverAddr = strings.ToLower(serverAddr)
	if !govalidator.IsDNSName(serverAddr) && !govalidator.IsIP(serverAddr) {
		return nil, ErrHandshakeHostname
	}

	serverPort := p.readUnsignedShort()
	if err != nil {
		return nil, ErrPacketParse
	}

	nextState, err := p.readVarint()
	if err != nil {
		return nil, ErrPacketParse
	}

	return &Handshake{
		Protoversion: protoVersion,
		ServerAddr:   serverAddr,
		ServerPort:   serverPort,
		NextState:    nextState,
	}, err
}

// Status

type StatusResponse struct {
	Version     VersionJSON     `json:"version"`
	Players     PlayersJSON     `json:"players"`
	Description json.RawMessage `json:"description"`
	Favicon     string          `json:"favicon"`
}

type VersionJSON struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type PlayersJSON struct {
	Max    int                `json:"max"`
	Online int                `json:"online"`
	Sample []PlayerSampleJSON `json:"sample"`
}

type PlayerSampleJSON struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

const StatusRequestID = 0x00

func (p *Packet) ReadStatusRequest() error {
	if p.ID != StatusRequestID {
		return ErrPacketID
	}
	return nil
}

const StatusResponseID = 0x00

func WriteStatusResponse(msg json.RawMessage) (*Packet, error) {
	length := writeVarint(uint32(len(msg)))

	data := append(length, msg...)

	return &Packet{
		ID:   StatusResponseID,
		Data: &data,
	}, nil
}

// Ping

const PingRequestID = 0x01

type Ping struct {
	Payload int64
}

func (p *Packet) ReadPingRequest() (*Ping, error) {
	if p.ID != PingRequestID {
		return nil, ErrPacketID
	}

	payload := p.readLong()

	return &Ping{
		Payload: payload,
	}, nil
}

const PingResponseID = 0x01

func WritePingResponse(ping *Ping) (*Packet, error) {
	num := ping.Payload
	data := []byte{
		byte(num >> 56), byte(num >> 48), byte(num >> 40), byte(num >> 32),
		byte(num >> 24), byte(num >> 16), byte(num >> 8), byte(num),
	}

	return &Packet{
		ID:   PingResponseID,
		Data: &data,
	}, nil
}
