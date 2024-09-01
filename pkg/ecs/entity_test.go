package ecs

import "testing"

func TestEntities(t *testing.T) {
	t.Run("spawns entity", func(t *testing.T) {
		em := newEntityManager()

		entity, err := em.create()
		if err != nil {
			t.Fail()
		}

		if entity > MaxEntities || entity < 0 {
			t.Fail()
		}
	})

	t.Run("destroys entity", func(t *testing.T) {
		em := newEntityManager()

		entityId, err := em.create()
		if err != nil {
			t.Fail()
		}

		err = em.destroy(entityId)
		if err != nil {
			t.Fail()
		}
	})

	t.Run("reuses entity id", func(t *testing.T) {
		em := newEntityManager()

		firstEntityId, err := em.create()
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		err = em.destroy(firstEntityId)
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		secondEntityId, err := em.create()
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		if firstEntityId != secondEntityId {
			t.Log(err)
			t.Fail()
		}
	})
}
