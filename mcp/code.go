package mcp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

// PLC Data communication code.
// This item is operating byte order.
type Code int

const (
	// Ascii code is normal mode.
	// Stored from upper byte to lower byte.
	Ascii Code = iota

	//　Binary code is approximately half the amount of communication data compared to communication using ASCII code
	// Stored from lower byte to upper byte.
	Binary
)

func (c Code) EncodeHex(s string) ([]byte, error) {
	if c == Ascii {
		return []byte(s), nil
	}

	decode, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	_ = binary.Write(buff, binary.LittleEndian, decode)
	return buff.Bytes(), nil
}
