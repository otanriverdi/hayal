package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcs(t *testing.T) {
	type transform struct {
		x int
		y int
	}

	t.Run("spawns entity", func(t *testing.T) {
		ecs := New()
		id, err := ecs.Spawn(5)
		assert.NoError(t, err)
		entityVal, ok := ecs.entityIndex.Load(id)
		entity := entityVal.(entityRef)
		assert.True(t, ok)
		assert.Equal(t, entity.row[0].(int), 5)
		assert.Len(t, ecs.archetypes, 1)
		assert.Len(t, ecs.archetypes[0].entities, 1)
		assert.Equal(t, ecs.archetypes[0].entities[0][0].(int), 5)
	})

	t.Run("destroys entity", func(t *testing.T) {
		ecs := New()
		id, err := ecs.Spawn(5)
		assert.NoError(t, err)
		err = ecs.Destroy(id)
		assert.NoError(t, err)
		_, ok := ecs.entityIndex.Load(id)
		assert.False(t, ok)
		assert.Len(t, ecs.archetypes, 1)
		assert.Len(t, ecs.archetypes[0].entities, 0)
	})

	t.Run("adds component", func(t *testing.T) {
		ecs := New()
		id, err := ecs.Spawn(5)
		assert.NoError(t, err)
		err = ecs.AddComponent(id, transform{x: 10, y: 5})
		assert.NoError(t, err)
		assert.Len(t, ecs.archetypes, 2)
		assert.Len(t, ecs.archetypes[0].entities, 0)
		assert.Len(t, ecs.archetypes[1].entities, 1)
		assert.Equal(t, ecs.archetypes[1].entities[0][0].(int), 5)
		assert.Equal(t, ecs.archetypes[1].entities[0][1].(transform).x, 10)
		assert.Equal(t, ecs.archetypes[1].entities[0][1].(transform).y, 5)
	})

	t.Run("removes component", func(t *testing.T) {
		ecs := New()
		id, err := ecs.Spawn(5)
		assert.NoError(t, err)
		err = ecs.AddComponent(id, transform{x: 10, y: 5})
		assert.NoError(t, err)
		err = ecs.RemoveComponent(id, transform{})
		assert.NoError(t, err)
		assert.Len(t, ecs.archetypes[0].entities, 1)
		assert.Len(t, ecs.archetypes[1].entities, 0)
		assert.Equal(t, ecs.archetypes[0].entities[0][0].(int), 5)
	})

	t.Run("queries components", func(t *testing.T) {
		ecs := New()
		_, err := ecs.Spawn(5)
		assert.NoError(t, err)
		_, err = ecs.Spawn(7)
		assert.NoError(t, err)
		iter, err := ecs.Query(5)
		assert.NoError(t, err)

		for res := range iter {
			_, err := GetComponent[int](&res)
			assert.NoError(t, err)
			err = SetComponent(&res, 3)
			assert.NoError(t, err)
		}
		assert.Equal(t, ecs.archetypes[0].entities[0][0].(int), 3)
	})
}
