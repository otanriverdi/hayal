package archetype

import (
	"testing"

	"github.com/otanriverdi/hayal/storage"
	"github.com/stretchr/testify/assert"
)

func TestArchetype(t *testing.T) {
	t.Run("creates archetype", func(t *testing.T) {
		a := NewArchetype(5, "foo")

		assert.Equal(t, len(a.storages), 2)
		assert.Equal(t, a.len, uint(0))
	})

	t.Run("hashes can be regenerated", func(t *testing.T) {
		a := NewArchetype(5, "foo")
		hash := Hash(5, "foo")
		unorderedHash := Hash("foo", 5)

		assert.Equal(t, a.hash, hash)
		assert.Equal(t, a.hash, unorderedHash)
	})

	t.Run("runs inclusion check", func(t *testing.T) {
		a := NewArchetype(5, "foo")

		res := DoesInclude(&a, "foo")
		assert.True(t, res)
		res = DoesInclude(&a, 1.5)
		assert.False(t, res)
	})

	t.Run("runs match check", func(t *testing.T) {
		a := NewArchetype(5, "foo")

		res := DoesMatch(&a, 5, "foo")
		assert.True(t, res)
		res = DoesMatch(&a, 5, "foo", 1.5)
		assert.False(t, res)
		res = DoesMatch(&a, 5)
		assert.False(t, res)
	})

	t.Run("can append data", func(t *testing.T) {
		a := NewArchetype(5, "foo")

		err := Append(&a, 3, "bar")
		assert.NoError(t, err)

		assert.Equal(t, a.len, uint(1))

		intSlice, err := storage.Slice[int](&a.storages[0])
		assert.NoError(t, err)
		assert.Equal(t, len(intSlice), 1)
		assert.Equal(t, intSlice[0], 3)

		strSlice, err := storage.Slice[string](&a.storages[1])
		assert.NoError(t, err)
		assert.Equal(t, len(strSlice), 1)
		assert.Equal(t, strSlice[0], "bar")
	})

	t.Run("prevents invalid append", func(t *testing.T) {
		a := NewArchetype(5, "foo")

		err := Append(&a, 3)
		assert.Error(t, err)

		err = Append(&a, "foo")
		assert.Error(t, err)

		err = Append(&a, "foo", 1.5)
		assert.Error(t, err)
	})

	t.Run("can remove data", func(t *testing.T) {
		a := NewArchetype(5, "foo")

		err := Append(&a, 3, "bar")
		assert.NoError(t, err)

		Remove(&a, 0)

		assert.Equal(t, a.len, uint(0))

		intSlice, err := storage.Slice[int](&a.storages[0])
		assert.NoError(t, err)
		assert.Equal(t, len(intSlice), 0)

		strSlice, err := storage.Slice[string](&a.storages[1])
		assert.NoError(t, err)
		assert.Equal(t, len(strSlice), 0)
	})
}
