package ecs

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	type transform struct {
		x int
		y int
	}

	t.Run("creates empty entity when it doesnt exist", func(t *testing.T) {
		cm := newComponentManager()
		db := newDatabase(&cm)

		var entityId uint16 = 1
		ref := db.ensureEntity(entityId)

		if ref.id != entityId {
			t.Log(ref)
			t.Fail()
		}
		if ref.archetypeId != 0 {
			t.Log(ref)
			t.Fail()
		}
		if ref.idx != 0 {
			t.Log(ref)
			t.Fail()
		}

		table, found := db.tables[0]
		if !found {
			t.Fail()
		}

		entity := table.entries[ref.idx]
		if entity.id != ref.id {
			t.Fail()
		}
		if len(entity.components) > 0 {
			t.Fail()
		}
	})

	t.Run("moves entity to new table after component insertion", func(t *testing.T) {
		cm := newComponentManager()
		db := newDatabase(&cm)

		var entityId uint16 = 1
		ref := db.ensureEntity(entityId)

		db.insertComponent(ref, transform{})

		if ref.archetypeId == 0 {
			t.Log(ref)
			t.Fail()
		}

		zeroTable, found := db.tables[0]
		if !found {
			t.Fail()
		}
		if len(zeroTable.entries) > 0 {
			t.Log(zeroTable)
			t.Fail()
		}

		actualTable, found := db.tables[ref.archetypeId]
		if !found {
			t.Fail()
		}
		entity := actualTable.entries[ref.idx]
		if entity.id != ref.id {
			t.Fail()
		}
		if len(entity.components) != 1 {
			t.Fail()
		}
	})

	t.Run("moves entity to new table after component deletion", func(t *testing.T) {
		cm := newComponentManager()
		db := newDatabase(&cm)

		var entityId uint16 = 1
		ref := db.ensureEntity(entityId)
		db.insertComponent(ref, transform{})
		oldTable, found := db.tables[ref.archetypeId]
		if !found {
			t.Fail()
		}

		db.deleteComponent(ref, transform{})

		table, found := db.tables[0]
		if !found {
			t.Fail()
		}
		entity := table.entries[ref.idx]
		if entity.id != ref.id {
			t.Fail()
		}
		if len(entity.components) > 0 {
			t.Fail()
		}
		if len(oldTable.entries) > 0 {
			t.Log(oldTable)
			t.Fail()
		}
	})

	t.Run("deletes the entity", func(t *testing.T) {
		cm := newComponentManager()
		db := newDatabase(&cm)

		var entityId uint16 = 1
		ref := db.ensureEntity(entityId)
		db.insertComponent(ref, transform{})
		db.deleteEntity(ref)
		oldTable, found := db.tables[ref.archetypeId]
		if !found {
			t.Fail()
		}
		if len(oldTable.entries) > 0 {
			t.Log(oldTable)
			t.Fail()
		}
		if len(db.refs) > 0 {
			t.Log(db)
			t.Fail()
		}
	})
}

func TestGetComponent(t *testing.T) {
	type transform struct {
		x int
		y int
	}

	t.Run("provides access to component data", func(t *testing.T) {
		cm := newComponentManager()
		db := newDatabase(&cm)

		var entityId uint16 = 1
		ref := db.ensureEntity(entityId)

		x := 10
		y := 15
		db.insertComponent(ref, transform{x, y})

		cmp, err := GetComponent[transform](&db, ref)
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		if cmp.x != x {
			t.Log(cmp)
			t.Fail()
		}
		if cmp.y != y {
			t.Log(cmp)
			t.Fail()
		}
	})

	t.Run("returns correct pointers to mutate data", func(t *testing.T) {
		cm := newComponentManager()
		db := newDatabase(&cm)

		var entityId uint16 = 1
		ref := db.ensureEntity(entityId)

		x := 10
		y := 15
		db.insertComponent(ref, transform{x, y})

		cmp, err := GetComponent[transform](&db, ref)
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		cmp.x = 25

		cmp, err = GetComponent[transform](&db, ref)
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		if cmp.x != 25 {
			t.Log(cmp)
			t.Fail()
		}
	})
}
