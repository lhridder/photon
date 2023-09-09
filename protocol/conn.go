package protocol

import (
	"fmt"
	"net"
)

const (
	Status uint = 1
	Login       = 2
)

type Conn struct {
	NetConn net.Conn
	Type    uint
}

func (c *Conn) readBytes(n uint) ([]byte, error) {
	buffer := make([]byte, n)
	count, err := c.NetConn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes: %s", err)
	}

	return buffer[:count], nil
}

func (c *Conn) readByte() (byte, error) {
	b, err := c.readBytes(1)
	return b[0], err
}

func (c *Conn) readPacketLength() (uint, error) {
	pos := 0
	var value int32
	for i := 0; i < 4; i++ {
		current, err := c.readByte()
		if err != nil {
			return 0, err
		}

		b := int32(current)
		value |= (b & 0x7F) << pos

		if current&0b10000000 == 0 {
			break
		}

		pos += 7
	}
	return uint(value), nil
}

// ReadPacket returns a packet and enforces a packet length limit measured in bytes
func (c *Conn) ReadPacket(limit uint) (*Packet, error) {
	length, err := c.readPacketLength()
	if err != nil || length < 0 {
		return nil, ErrPacketLengthInvalid
	}

	if length > limit {
		return nil, ErrPacketLengthLimit
	}

	bytes, err := c.readBytes(length)
	if err != nil {
		//TODO check conn closed or time exceeded
		return nil, ErrPacketRead
	}

	data := bytes[1:]

	pk := Packet{
		ID:   bytes[0],
		Data: &data,
	}

	return &pk, nil
}

// WritePacket writes a *Packet to *Conn
func (c *Conn) WritePacket(pk *Packet) error {
	length := len(*pk.Data) + 1
	packetsize := writeVarint(uint32(length))
	data := append(packetsize, pk.ID)
	data = append(data, *pk.Data...)

	_, err := c.NetConn.Write(data)
	if err != nil {
		return err
	}

	return nil
}
