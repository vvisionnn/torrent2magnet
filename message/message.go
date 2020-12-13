package message

import (
	"encoding/binary"
	"io"
)

type messageID uint8

// Message store the message id and payload
type Message struct {
	ID 		messageID
	Payload	[]byte
}

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

// Read read a message from a stream, return nil on keep-alive message
func Read(r io.Reader) (*Message, error) {
	// the first 32 bit contains the length of this message
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// keep-alive message
	if length == 0 { return nil, nil }

	// get message buf
	messageBuf := make([]byte, length)
	if _, err := io.ReadFull(r, messageBuf); err != nil {
		return nil, err
	}

	m := Message{
		// the first 8 bit of message buf is message id
		ID: messageID(messageBuf[0]),
		Payload: messageBuf[1:],	// and last if message payload
	}
	return &m, nil
}

// Serialize serializes a message into a buffer of form
// <length><message id><payload>
// Interprets `nil` as a keep-alive message
func (m *Message) Serialize () []byte {
	if m == nil { return make([]byte, 4) }

	length := uint32(len(m.Payload) + 1) // +1 for id
	buf := make([]byte, length + 4)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

// FormatRequest creates a REQUEST message
func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: MsgRequest, Payload: payload}
}

// FormatHave creates a HAVE message
func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{ID: MsgHave, Payload: payload}
}