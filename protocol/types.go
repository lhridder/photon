package protocol

func (p *Packet) readVarint() (int32, error) {
	data := *p.Data
	pos := 0
	var value int32
	for i := 0; i < 4; i++ {
		current := data[i]
		b := int32(current)
		value |= (b & 0x7F) << pos

		if current&0b10000000 == 0 {
			*p.Data = data[i+1:]
			break
		}

		pos += 7
	}

	return value, nil
}

func writeVarint(num uint32) []byte {
	var value []byte
	for {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		value = append(value, byte(b))
		if num == 0 {
			break
		}
	}

	return value
}

func (p *Packet) readString() (string, error) {
	stringLength, err := p.readVarint()
	if err != nil {
		return "", err
	}

	data := *p.Data
	text := string(data[:stringLength])

	*p.Data = data[stringLength:]

	return text, nil
}

func (p *Packet) readUnsignedShort() uint16 {
	data := *p.Data

	num := uint16(data[0])<<8 | uint16(data[1])

	*p.Data = data[2:]

	return num
}

func (p *Packet) readLong() int64 {
	data := *p.Data

	num := int64(data[0])<<56 | int64(data[1])<<48 | int64(data[2])<<40 | int64(data[3])<<32 |
		int64(data[4])<<24 | int64(data[5])<<16 | int64(data[6])<<8 | int64(data[7])

	*p.Data = data[8:]

	return num
}
