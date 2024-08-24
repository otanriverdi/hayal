package ecs

import (
	"testing"
)

func TestEntities(t *testing.T) {
	t.Run("creates entity", func(t *testing.T) {
		ecs := New()

		entityId, err := ecs.CreateEntity()
		if err != nil {
			t.Fail()
		}

		if entityId > MaxEntitites || entityId < 0 {
			t.Fail()
		}
	})

	t.Run("destroys entity", func(t *testing.T) {
		ecs := New()

		entityId, err := ecs.CreateEntity()
		if err != nil {
			t.Fail()
		}

		err = ecs.DestroyEntity(entityId)
		if err != nil {
			t.Fail()
		}
	})

	t.Run("reuses entity id", func(t *testing.T) {
		ecs := New()

		firstEntityId, err := ecs.CreateEntity()
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		err = ecs.DestroyEntity(firstEntityId)
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		secondEntityId, err := ecs.CreateEntity()
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

func TestComponents(t *testing.T) {
	type transform struct {
		x int
	}

	ecs := New()

	entityId, err := ecs.CreateEntity()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Run("adds component by passing a value", func(t *testing.T) {
		cmp := transform{}
		err = ecs.AddComponent(entityId, &cmp)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	})

	t.Run("adds component by passing a pointer", func(t *testing.T) {
		newEntityId, err := ecs.CreateEntity()
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		cmp := transform{}
		err = ecs.AddComponent(newEntityId, cmp)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	})

	t.Run("adds component by passing anonymous struct", func(t *testing.T) {
		newEntityId, err := ecs.CreateEntity()
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		err = ecs.AddComponent(newEntityId, transform{})
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	})

	t.Run("validates component size", func(t *testing.T) {
		type oversized struct {
			arr [65]byte
		}
		cmp := oversized{}
		err = ecs.AddComponent(entityId, cmp)
		if err == nil {
			t.Fail()
		}
		if err != ErrComponentSizeTooBig {
			t.Log(err)
			t.Fail()
		}
	})

	t.Run("doesnt allow duplicate components", func(t *testing.T) {
		cmp := transform{}
		err = ecs.AddComponent(entityId, cmp)
		if err == nil {
			t.Fail()
		}

		if err != ErrComponentAlreadySlotted {
			t.Log(err)
			t.Fail()
		}
	})

	t.Run("removes component", func(t *testing.T) {
		err = ecs.RemoveComponent(entityId, transform{})
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	})
}

func TestGetComponent(t *testing.T) {
	type transform struct {
		x int
		y int
	}

	ecs := New()

	entityId, err := ecs.CreateEntity()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Run("provides access to component data", func(t *testing.T) {
		x := 10
		y := 15
		tr := transform{x, y}
		err = ecs.AddComponent(entityId, tr)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		cmp, err := GetEcsComponent(&ecs, entityId, transform{})
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if cmp.x != x {
			t.Log(cmp.x)
			t.Fail()
		}
		if cmp.y != y {
			t.Log(cmp.y)
			t.Fail()
		}
	})
}

func TestQuery(t *testing.T) {
	type transform struct {
		x int
		y int
	}

	type name struct {
		name string
	}

	type velocity struct {
		velocity int
	}

	t.Run("returns correct entity for no match", func(t *testing.T) {
		ecs := New()
		entityId, err := ecs.CreateEntity()
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		ecs.AddComponent(entityId, transform{})
		ecs.AddComponent(entityId, velocity{})

		secondEntityId, err := ecs.CreateEntity()
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		ecs.AddComponent(secondEntityId, name{})

		results, err := ecs.Query(transform{}, velocity{})
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		if len(results) != 1 {
			t.Log(len(results))
			t.Fail()
		}

		if results[0] != entityId {
			t.Log(results[0])
			t.Fail()
		}
	})

	t.Run("returns correct entity for partial match", func(t *testing.T) {
		ecs := New()
		entityId, err := ecs.CreateEntity()
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		ecs.AddComponent(entityId, transform{})
		ecs.AddComponent(entityId, velocity{})

		secondEntityId, err := ecs.CreateEntity()
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		ecs.AddComponent(secondEntityId, name{})
		ecs.AddComponent(secondEntityId, transform{})

		results, err := ecs.Query(transform{}, velocity{})
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		if len(results) != 1 {
			t.Log(len(results))
			t.Fail()
		}

		if results[0] != entityId {
			t.Log(results[0])
			t.Fail()
		}
	})
}
