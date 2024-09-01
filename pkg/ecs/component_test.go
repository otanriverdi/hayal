package ecs

import (
	"testing"
)

func TestComponents(t *testing.T) {
	type transform struct {
		x int
	}

	t.Run("registers component by passing a value", func(t *testing.T) {
		cm := newComponentManager()
		cmp := transform{}
		expectedCmpId := cm.ensureId(&cmp)
		receivedCmpId, err := cm.getId(&cmp)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if expectedCmpId != receivedCmpId {
			t.Fail()
		}
	})

	t.Run("registers component by passing a pointer", func(t *testing.T) {
		cm := newComponentManager()
		cmp := transform{}
		expectedCmpId := cm.ensureId(cmp)
		receivedCmpId, err := cm.getId(&cmp)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if expectedCmpId != receivedCmpId {
			t.Fail()
		}
	})

	t.Run("registers component by passing anonymous struct", func(t *testing.T) {
		cm := newComponentManager()
		expectedCmpId := cm.ensureId(transform{})
		receivedCmpId, err := cm.getId(transform{})
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if expectedCmpId != receivedCmpId {
			t.Fail()
		}
	})

	t.Run("returns not found for get operation when component doesnt exist", func(t *testing.T) {
		cm := newComponentManager()
		_, err := cm.getId(transform{})
		if err == nil {
			t.Fail()
		}
		if err != ErrComponentNotFound {
			t.Log(err)
			t.Fail()
		}
	})
}
