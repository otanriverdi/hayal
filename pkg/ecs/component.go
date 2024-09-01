package ecs

import (
	"errors"
	"reflect"
)

var (
	ErrComponentNotFound   = errors.New("component not found")
	ErrComponentNilPointer = errors.New("component is nil pointer")
)

// InvalidComponentId is the NOOP component id returned when related operations throw
const InvalidComponentId = ^componentId(0)

// ComponentId is a unique identifier for a single component.
type componentId = uint16

type componentManager struct {
	count    componentId
	registry map[reflect.Type]componentId
}

func newComponentManager() componentManager {
	return componentManager{
		count:    0,
		registry: make(map[reflect.Type]componentId),
	}
}

func (cm *componentManager) ensureId(cmp any) componentId {
	cmpType := cm.deriveType(cmp)
	cmpId, found := cm.registry[cmpType]
	if !found {
		cmpId = cm.count
		cm.registry[cmpType] = cmpId
		cm.count++
	}
	return cmpId
}

func (cm *componentManager) getId(cmp any) (componentId, error) {
	cmpType := cm.deriveType(cmp)
	cmpId, found := cm.registry[cmpType]
	if !found {
		return InvalidComponentId, ErrComponentNotFound
	}
	return cmpId, nil
}

func (cm *componentManager) deriveType(cmp any) reflect.Type {
	cmpType := reflect.TypeOf(cmp)
	if cmpType.Kind() == reflect.Ptr {
		cmpType = cmpType.Elem()
	}
	return cmpType
}
