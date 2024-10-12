package ecs

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
)

type componentId = uint32

var (
	cmpIdInc componentId
	cmpIdMap sync.Map
)

const MAX_COMPONENTS uint8 = 128

type bitmap = [MAX_COMPONENTS / 64]uint64

type archetype struct {
	bitmap     bitmap
	entities   [][]any
	cmpIndices map[componentId]int
	mu sync.Mutex
}

type entity = uint64

var entityIdInc entity

type entityRef struct {
	idx    int
	bitmap bitmap
	row    []any
}

type ECS struct {
	archetypes     []archetype
	archetypeIndex sync.Map 
	entityIndex    sync.Map 
	// For archetypes array
	mu sync.RWMutex
}

func New() ECS {
	return ECS{
		archetypes:     make([]archetype, 0),
	}
}

func (ecs *ECS) Spawn(cmp any) (entity, error) {
	id := atomic.AddUint64(&entityIdInc, 1)
	cmpId, err := getCmpId(cmp)
	if err != nil {
		return 0, err
	}
	bitmap := buildBitmap(cmpId)
	ref, err := ecs.insertRow(id, bitmap, cmp)
	if err != nil {
		return 0, err
	}
	ecs.entityIndex.Store(id, ref);
	return id, nil
}

func (ecs *ECS) Destroy(entity entity) error {
	refVal, ok := ecs.entityIndex.LoadAndDelete(entity)
	if !ok {
		return errors.New("Entity not found")
	}
	ref := refVal.(entityRef)
	ecs.deleteRow(ref.bitmap, ref.idx)
	return nil
}

func (ecs *ECS) AddComponent(entity entity, cmp any) error {
	refVal, ok := ecs.entityIndex.Load(entity)
	if !ok {
		return errors.New("Entity not found")
	}
	ref := refVal.(entityRef)
	cmpId, err := getCmpId(cmp)
	if err != nil {
		return err
	}
	bitmap := setBitmap(ref.bitmap, cmpId)
	cmps := append(ref.row, cmp)
	newRef, err := ecs.insertRow(entity, bitmap, cmps...)
	if err != nil {
		return err
	}
	ecs.entityIndex.Store(entity, newRef);
	ecs.deleteRow(ref.bitmap, ref.idx)
	return nil
}

func (ecs *ECS) RemoveComponent(entity entity, cmp any) error {
	refVal, ok := ecs.entityIndex.Load(entity)
	if !ok {
		return errors.New("Entity not found")
	}
	ref := refVal.(entityRef)
	cmpId, err := getCmpId(cmp)
	if err != nil {
		return err
	}
	bitmap := clearBitmap(ref.bitmap, cmpId)
	cmpIdx := len(ref.row)
	for idx, cmp := range ref.row {
		rowCmpId, err := getCmpId(cmp)
		if err != nil {
			return err
		}
		if rowCmpId == cmpId {
			cmpIdx = idx
			break
		}
	}
	cmps := append(ref.row[0:cmpIdx], ref.row[cmpIdx+1:len(ref.row)]...)
	newRef, err := ecs.insertRow(entity, bitmap, cmps...)
	if err != nil {
		return err
	}
	ecs.entityIndex.Store(entity, newRef);
	ecs.deleteRow(ref.bitmap, ref.idx)
	return nil
}

func (ecs *ECS) Query(cmps ...any) (func(yield func(QueryResult) bool), error) {
	cmpIds := make([]uint32, len(cmps))
	for i, cmp := range cmps {
		cmpId, err := getCmpId(cmp)
		if err != nil {
			return nil, err
		}
		cmpIds[i] = cmpId
	}
	queryBitmap := buildBitmap(cmpIds...)
	return func(yield func(QueryResult) bool) {
		for i := range ecs.archetypes {
			a := &ecs.archetypes[i]
			if !bitmapIsSubset(queryBitmap, a.bitmap) {
				continue
			}
			for idx := range a.entities {
				components := a.entities[idx]
				qr := QueryResult{
					components:   components,
					cmpIndices:   a.cmpIndices,
					archetype:    a,
					archetypeIdx: idx,
				}
				if !yield(qr) {
					return
				}
			}
		}
	}, nil
}

func (ecs *ECS) deleteRow(bitmap bitmap, idx int) {
	a := ecs.ensureArchetype(bitmap)
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entities = append(a.entities[0:idx], a.entities[idx+1:len(a.entities)]...)
}

func (ecs *ECS) insertRow(entity entity, bitmap bitmap, cmps ...any) (entityRef, error) {
	a := ecs.ensureArchetype(bitmap)
	a.mu.Lock()
	defer a.mu.Unlock()
	row := make([]any, len(cmps))
	for _, cmp := range cmps {
		cmpId, err := getCmpId(cmp)
		if err != nil {
			return entityRef{}, err
		}
		idx := a.cmpIndices[cmpId]
		row[idx] = cmp
	}
	a.entities = append(a.entities, row)
	idx := len(a.entities) - 1
	return entityRef{row: a.entities[idx], idx: idx, bitmap: bitmap}, nil
}

func (ecs *ECS) ensureArchetype(bitmap bitmap) *archetype {
	idxVal, ok := ecs.archetypeIndex.Load(bitmap)
	if !ok {
		ecs.mu.Lock()
		defer ecs.mu.Unlock()
		cmpIds := extractBitmapCmps(bitmap)
		cmpIndices := make(map[componentId]int)
		for idx, cmpId := range cmpIds {
			cmpIndices[cmpId] = idx
		}
		ecs.archetypes = append(ecs.archetypes, archetype{bitmap: bitmap, entities: make([][]any, 0), cmpIndices: cmpIndices})
		idx := len(ecs.archetypes) - 1
		idxVal = idx
		ecs.archetypeIndex.Store(bitmap, idx)
	}
	idx := idxVal.(int)
	return &ecs.archetypes[idx]
}

func getCmpId(cmp any) (componentId, error) {
	cmpType := reflect.TypeOf(cmp)
	if id, ok := cmpIdMap.Load(cmpType); ok {
		return id.(componentId), nil
	}
	if cmpIdInc >= uint32(MAX_COMPONENTS) {
		return 0, errors.New("Max number of components")
	}
	id := atomic.AddUint32(&cmpIdInc, 1)
	cmpIdMap.Store(cmpType, id)
	return id, nil
}

type QueryResult struct {
	components []any
	cmpIndices    map[componentId]int
	archetype *archetype
	archetypeIdx int
}

func GetComponent[C any](qr *QueryResult) (C, error) {
	var zero C
	cmpId, err := getCmpId(zero)
	if err != nil {
		return zero, err
	}
	idx, ok := qr.cmpIndices[cmpId]
	if !ok {
		return zero, errors.New("Component in type param does not exist in this query result")
	}
	if _, ok := qr.components[idx].(C); !ok {
		panic("Improperly mapped cmp indices")
	}
	component, ok := qr.components[idx].(C)
	if !ok {
		return zero, errors.New("Component type does not match expected type")
	}
	return component, nil
}

func SetComponent(qr *QueryResult, cmp any) error {
	cmpId, err := getCmpId(cmp)
	if err != nil {
		return err
	}
	idx, ok := qr.cmpIndices[cmpId]
	if !ok {
		return errors.New("Component does not exist in this query result")
	}
	qr.archetype.entities[qr.archetypeIdx][idx] = cmp
	return nil
}
