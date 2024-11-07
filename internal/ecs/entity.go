package ecs

import (
	"sync/atomic"
)

type Entity uint64

var (
	currentEntityId uint64
)

func NewEntity() Entity {
	return Entity(atomic.AddUint64(&currentEntityId, 1))
}
