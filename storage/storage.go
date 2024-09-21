package storage

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var cmpTypeCache sync.Map

// Storage is a type erased slice for a single component.
type Storage struct {
	slice reflect.Value
}

func New(cmp any) Storage {
	cmpType := reflect.TypeOf(cmp)
	sliceType := reflect.SliceOf(cmpType)
	slice := reflect.MakeSlice(sliceType, 0, 10)
	return Storage{slice}
}

// Append appends the component to the storage providing the types match.
func Append(s *Storage, cmp any) error {
	cmpType := Type(cmp)
	if cmpType != s.slice.Type().Elem() {
		return errors.New(fmt.Sprintf("Unexpected component type. Expected: %s, Got: %s", s.slice.Type().String(), cmpType.String()))
	}
	value := reflect.ValueOf(cmp)
	s.slice = reflect.Append(s.slice, value)
	return nil
}

// Remove removes the component from the storage at the provided index.
func Remove(s *Storage, idx uint) {
	i := int(idx)
	if i + 1 > s.slice.Len() {
			return
	}
	s.slice = reflect.AppendSlice(s.slice.Slice(0, i), s.slice.Slice(i+1, s.slice.Len()))
}

// Slice returns the storage slice casted to the type parameter provided the types match.
func Slice[T any](s *Storage) ([]T, error) {
	slice, ok := s.slice.Interface().([]T)
	if !ok {
		return nil, errors.New("Type parameter does not match storage type")
	}
	return slice, nil
}

func At[T any](s *Storage, idx uint) (*T, error) {
	slice, error := Slice[T](s)
	if error != nil {
		return nil, error
	}
	return &slice[idx], nil
}

// IsType is a predicate to determine if the type of component matches storage type.
func IsType(s *Storage, cmp any) bool {
	cmpType := Type(cmp)
	return cmpType == s.slice.Type().Elem()
}

// Type returns the type of the component
func Type(cmp any) reflect.Type {
	if id, ok := cmpTypeCache.Load(cmp); ok {
		return id.(reflect.Type)
	}
	t := reflect.TypeOf(cmp)
	cmpTypeCache.Store(cmp, t)
	return t
}

