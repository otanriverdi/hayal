package archetype

import (
	"encoding/binary"
	"errors"
	"hash/fnv"
	"sort"

	"github.com/otanriverdi/hayal/storage"
)

// Archetype is the storage structure for entities that shares the same group of components.
type Archetype struct {
	len      uint
	storages []storage.Storage
	hash     uint64
}

func NewArchetype(cmps ...any) Archetype {
	storages := make([]storage.Storage, len(cmps), len(cmps))
	for idx, cmp := range cmps {
		storages[idx] = storage.New(cmp)
	}
	hash := Hash(cmps...)
	return Archetype{0, storages, hash}
}

// DoesInclude checks if the archetype includes the passed in component types.
func DoesInclude(a *Archetype, cmps ...any) bool {
main:
	for _, cmp := range cmps {
		for _, s := range a.storages {
			if storage.IsType(&s, cmp) {
				continue main
			}
		}
		return false
	}
	return true
}

// DoesMatch checks if the archetype matches the passed in component types exactly.
func DoesMatch(a *Archetype, cmps ...any) bool {
	if len(cmps) != len(a.storages) {
		return false
	}
	return DoesInclude(a, cmps...)
}

// Append appends the entity to archetype provided they match the archetype.
func Append(a *Archetype, cmps ...any) error {
	if !DoesMatch(a, cmps...) {
		return errors.New("Component types does not match the archetype")
	}
main:
	for _, cmp := range cmps {
		for idx := range a.storages {
			s := &a.storages[idx]
			if storage.IsType(s, cmp) {
				storage.Append(s, cmp)
				continue main
			}
		}
	}
	a.len++
	return nil
}

// Remove removes the entity at the provided index
func Remove(a *Archetype, idx uint) {
	if idx+1 > a.len {
		return
	}
	for i := range a.storages {
		s := &a.storages[i]
		storage.Remove(s, idx)
	}
	a.len--
}

// Hash creates a hash for the component types that is guaranteed to be unique for program execution. The resulting
// hash of this function should not be used for serializing outside of the program as in the next run the hash
// might be different.
func Hash(cmps ...any) uint64 {
	var typeHashes []uint64
	for _, cmp := range cmps {
		t := storage.Type(cmp)
		h := fnv.New64a()
		h.Write([]byte(t.String()))
		typeHashes = append(typeHashes, h.Sum64())
	}
	sort.Slice(typeHashes, func(i, j int) bool { return typeHashes[i] < typeHashes[j] })
	finalHash := fnv.New64a()
	for _, th := range typeHashes {
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], th)
		finalHash.Write(buf[:])
	}
	return finalHash.Sum64()
}
