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

	var cmpId uint16 = 1
	ecs := New()

	entityId, err := ecs.CreateEntity()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Run("adds component by passing a value", func(t *testing.T) {
		cmp := transform{}
		err = ecs.AddComponent(entityId, cmpId, &cmp)
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
		err = ecs.AddComponent(newEntityId, cmpId, cmp)
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
		err = ecs.AddComponent(newEntityId, cmpId, transform{})
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	})

	t.Run("validates component size", func(t *testing.T) {
		var oversizeCmpId uint16 = 2
		type oversized struct {
			arr [65]byte
		}
		cmp := oversized{}
		err = ecs.AddComponent(entityId, oversizeCmpId, cmp)
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
		err = ecs.AddComponent(entityId, cmpId, cmp)
		if err == nil {
			t.Fail()
		}

		if err != ErrComponentAlreadySlotted {
			t.Log(err)
			t.Fail()
		}
	})

	t.Run("removes component", func(t *testing.T) {
		err = ecs.RemoveComponent(entityId, cmpId)
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

	var cmpId uint16 = 1
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
		err = ecs.AddComponent(entityId, cmpId, tr)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		cmp, err := GetEcsComponent[transform](&ecs, entityId, cmpId)
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
