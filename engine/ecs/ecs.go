// Package ecs implements an (E)ntity (C)omponent (S)ystem.
//
// Entities are all beings that exist in a game. Groups of Components can be assigned to entities that would define their
// properties. System can process groups of entities that share the same components.
package ecs

import (
	"errors"
	"reflect"
	"unsafe"
)

const (
	MaxEntitites     = 2048
	MaxComponents    = 64
	MaxComponentSize = 64
)

type (
	// EntityId is a unique identifier for a single entity.
	EntityId = uint32
	// ComponentId is a unique identifier for a single component.
	ComponentId = uint16
	// TypeId is a bitmask identifier for groups of components as types.
	TypeId = uint64
	// System is the function type that can operate on entities of certain type.
	System = func(ecs *Ecs)
)

var (
	ErrComponentAlreadySlotted = errors.New("component slot is already in use")
	ErrEntityOutOfBounds       = errors.New("entity out of bounds")
	ErrComponentSizeTooBig     = errors.New("component size exceeds limit")
	ErrComponentNotSlotted     = errors.New("component not slotted")
	ErrEntityNotFound          = errors.New("entity not found")
	ErrMaxCapacity             = errors.New("max entity capacity")
)

type Ecs struct {
	freedCount uintptr
	freedIds   [MaxEntitites]EntityId

	activeCount    uintptr
	activeEntities [MaxEntitites]EntityId

	entityCount uintptr
	entityTypes [MaxEntitites]TypeId

	storage [MaxComponents][MaxComponentSize * MaxEntitites]uint8
}

func setBit(mask *TypeId, bit ComponentId) {
	*mask |= (1 << bit)
}

func clearBit(mask *TypeId, bit ComponentId) {
	*mask &= ^(1 << bit)
}

func cmpBit(mask TypeId, bit ComponentId) bool {
	return (mask & (1 << bit)) != 0
}

func (ecs *Ecs) CreateEntity() (EntityId, error) {
	var entityId EntityId

	if ecs.freedCount > 0 {
		ecs.freedCount--
		entityId = ecs.freedIds[ecs.freedCount]
	} else if ecs.entityCount < MaxEntitites {
		entityId = EntityId(ecs.entityCount)
		ecs.entityCount++
	} else {
		// EntityId of max uint32 is used as an invalid ID since it can't be assigned due to above condition
		return ^EntityId(0), ErrMaxCapacity
	}

	return entityId, nil
}

func (ecs *Ecs) DestroyEntity(entityId EntityId) error {
	ecs.entityTypes[entityId] = 0
	ecs.freedIds[ecs.freedCount] = entityId
	ecs.freedCount++

	for i := uintptr(0); i < ecs.activeCount; i++ {
		// Replace the entity with the last entity in the list
		if ecs.activeEntities[i] == entityId {
			ecs.activeEntities[i] = ecs.activeEntities[ecs.activeCount-1]
			ecs.activeCount--
			return nil
		}
	}

	return ErrEntityNotFound
}

func (ecs *Ecs) AddComponent(entityId EntityId, cmpId ComponentId, data *any) error {
	if cmpBit(ecs.entityTypes[entityId], cmpId) {
		return ErrComponentAlreadySlotted
	}

	if entityId >= MaxEntitites {
		return ErrEntityOutOfBounds
	}

	cmpSize := reflect.TypeOf(data).Elem().Size()
	if cmpSize > MaxComponentSize {
		return ErrComponentSizeTooBig
	}

	setBit(&ecs.entityTypes[entityId], cmpId)

	// Go does not allow converting generic types into byte arrays so we cast it into an unsafe pointer
	dataPtr := unsafe.Pointer(data)
	// We get a slice of bytes that represents the component data
	dataSlice := (*[MaxComponentSize]byte)(dataPtr)[:cmpSize:cmpSize]
	// We calculate the offset of this component data that will be sloted into the storage buffer
	offset := entityId * MaxComponentSize
	// Get the slice of the storage buffer that is the slot for this component data
	storageSlice := ecs.storage[cmpId][offset : offset+uint32(cmpSize)]
	// Copy component data
	copy(storageSlice, dataSlice)

	return nil
}

func (ecs *Ecs) RemoveComponent(entityId EntityId, cmpId ComponentId) error {
	if !cmpBit(ecs.entityTypes[entityId], cmpId) {
		return ErrComponentNotSlotted
	}

	if entityId >= MaxEntitites {
		return ErrEntityOutOfBounds
	}

	clearBit(&ecs.entityTypes[entityId], cmpId)
	return nil
}
