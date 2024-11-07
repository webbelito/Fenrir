// internal/network/message_id_generator.go

package network

import (
	"sync/atomic"
)

type MessageIDGenerator struct {
	currentID uint64
}

func NewMessageIDGenerator() *MessageIDGenerator {
	return &MessageIDGenerator{
		currentID: 0,
	}
}

func (m *MessageIDGenerator) NextID() uint64 {
	return atomic.AddUint64(&m.currentID, 1)
}
