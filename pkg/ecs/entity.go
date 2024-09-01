package ecs

import "errors"

var (
	ErrEntityNotFound    = errors.New("entity not found")
	ErrMaxEntityCapacity = errors.New("max entity capacity")
)

// MaxEntities represents the number of entities that can exist in the ecs
const MaxEntities = ^entityId(0)

type entityId = uint16

type entityManager struct {
	freed  []entityId
	active []entityId
	size   entityId
}

func newEntityManager() entityManager {
	return entityManager{}
}

func (em *entityManager) create() (entityId, error) {
	var entity entityId

	if len(em.freed) > 0 {
		entity = em.freed[len(em.freed)-1]
	} else if em.size < MaxEntities {
		entity = entityId(em.size)
		em.size++
	} else {
		// EntityId of max uint32 is used as an invalid ID since it can't be assigned due to above condition
		return MaxEntities, ErrMaxEntityCapacity
	}

	em.active = append(em.active, entity)

	return entity, nil
}

func (em *entityManager) destroy(entity entityId) error {
	em.freed = append(em.freed, entity)

	for idx, e := range em.active {
		if e == entity {
			em.active[idx] = em.active[len(em.active)-1]
			em.active = em.active[:len(em.active)-1]
			return nil
		}
	}

	return ErrEntityNotFound
}
