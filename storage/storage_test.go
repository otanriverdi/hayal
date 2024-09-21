package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	type transform struct {
		x int
		y int
	}

	t.Run("creates new storage from orphan value", func(t *testing.T) {
		s := New(transform{})

		assert.Equal(t, s.slice.Type().String(), "[]storage.transform")
		assert.Equal(t, s.slice.Len(), 0)
	})

	t.Run("creates new storage from owned value", func(t *testing.T) {
		cmp := transform{}
		s := New(cmp)

		assert.Equal(t, s.slice.Type().String(), "[]storage.transform")
		assert.Equal(t, s.slice.Len(), 0)
	})

	t.Run("creates new storage from primitive", func(t *testing.T) {
		s := New(uint8(5))

		assert.Equal(t, s.slice.Type().String(), "[]uint8")
		assert.Equal(t, s.slice.Len(), 0)
	})

	t.Run("can append to storage and retrieve slice", func(t *testing.T) {
		s := New(transform{})
		cmp := transform{x: 5, y: 10}
		err := Append(&s, cmp)

		assert.NoError(t, err)

		slice, err := Slice[transform](&s)

		assert.NoError(t, err)
		assert.Equal(t, slice[0].x, 5)
		assert.Equal(t, slice[0].y, 10)
	})

	t.Run("can retrieve single item storage", func(t *testing.T) {
		s := New(transform{})
		cmp := transform{x: 5, y: 10}
		err := Append(&s, cmp)

		assert.NoError(t, err)

		tr, err := At[transform](&s, 0)

		assert.NoError(t, err)
		assert.Equal(t, tr.x, 5)
		assert.Equal(t, tr.y, 10)
	})

	t.Run("rejects appends from incompatible values", func(t *testing.T) {
		s := New(transform{})
		err := Append(&s, 5)

		assert.Error(t, err)
	})

	t.Run("can remove value from storage", func(t *testing.T) {
		s := New(5)
		err := Append(&s, 5)

		assert.NoError(t, err)

		Remove(&s, 0)

		assert.NoError(t, err)
		assert.Equal(t, s.slice.Len(), 0)
	})

	t.Run("compare type", func(t *testing.T) {
		s := New(5)
		err := Append(&s, 5)

		assert.NoError(t, err)

		res := IsType(&s, 5)
		assert.True(t, res)
		res = IsType(&s, "foo")
		assert.False(t, res)
	})
}
