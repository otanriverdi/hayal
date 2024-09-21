package entity

import "sync/atomic"

var idInc uint64

// Entity is the representation of a single game object as a unique identifier
type Entity = uint64

func NewEntity() Entity {
	return atomic.AddUint64(&idInc, 1)
}
