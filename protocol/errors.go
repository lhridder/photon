package protocol

import "errors"

var (
	ErrPacketLengthInvalid = errors.New("packet length invalid")
	ErrPacketLengthLimit   = errors.New("packet length limit reached")
	ErrPacketRead          = errors.New("packet read failed")
	ErrPacketID            = errors.New("packet id incorrect")
	ErrPacketParse         = errors.New("packet parse failed")
	ErrHandshakeHostname   = errors.New("handshake hostname incorrect")
	ErrJsonMarshal         = errors.New("json marshal failed")
)
