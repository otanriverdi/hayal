package ecs

import "reflect"

// MaxComponents constant represents the number of maximum components the archetype storage can store
const MaxComponents = 16

// Query returns an iterator of all entities containing the components in the arguments list
func (db *database) Query(cmps ...any) func(yield func(*entity) bool) {
	return func(yield func(ref *entity) bool) {
		for archetypeId, table := range db.tables {
			matches := true
			for _, cmp := range cmps {
				cmpId := db.cm.ensureId(cmp)
				if !cmpBit(archetypeId, cmpId) {
					matches = false
					break
				}
			}
			if matches {
				for _, entry := range table.entries {
					if !yield(entry) {
						return
					}
				}
			}
		}
	}
}

// GetComponent extracts the component data from storage
func GetComponent[T any](db *database, e *entity) (*T, error) {
	cmpId, err := db.cm.getId((*T)(nil))
	if err != nil {
		return nil, err
	}
	if !cmpBit(e.archetypeId, cmpId) {
		return nil, ErrComponentNotFound
	}
	data := db.tables[e.archetypeId].entries[e.idx].components[cmpId]
	typedData, ok := data.(*T)
	if !ok {
		panic("Mismatching component type retrieved from storage")
	}
	return typedData, nil
}

type archetypeId = uint64

type table struct {
	archetypeId archetypeId
	entries     []*entity
}

type entity struct {
	id          entityId
	archetypeId archetypeId
	idx         int
	components  map[componentId]any
}

type database struct {
	tables map[archetypeId]*table
	refs   map[entityId]entity // TODO: how about making this an array?
	cm     *componentManager
}

func newDatabase(cm *componentManager) database {
	return database{
		tables: make(map[archetypeId]*table),
		refs:   make(map[entityId]entity),
		cm:     cm,
	}
}

func (db *database) ensureEntity(entityId entityId) *entity {
	ref, found := db.refs[entityId]
	if found {
		return &ref
	}
	ref = entity{
		id:          entityId,
		archetypeId: 0,
		idx:         0,
		components:  make(map[componentId]any),
	}
	t := db.ensureTable(ref.archetypeId)
	ref.idx = len(t.entries)
	t.entries = append(t.entries, &ref)
	return &ref
}

func (db *database) insertComponent(e *entity, cmp any) {
	cmpId, cmp := db.parseComponent(cmp)
	e.components[cmpId] = cmp
	if !cmpBit(e.archetypeId, cmpId) {
		db.tables[e.archetypeId].entries = append(
			db.tables[e.archetypeId].entries[:e.idx],
			db.tables[e.archetypeId].entries[e.idx+1:]...,
		)
		setBit(&e.archetypeId, cmpId)
		t := db.ensureTable(e.archetypeId)
		e.idx = len(t.entries)
		t.entries = append(t.entries, e)
	}
}

func (db *database) deleteComponent(e *entity, cmp any) error {
	db.tables[e.archetypeId].entries = append(
		db.tables[e.archetypeId].entries[:e.idx],
		db.tables[e.archetypeId].entries[e.idx+1:]...,
	)
	cmpId, err := db.cm.getId(cmp)
	if err != nil {
		return err
	}
	delete(e.components, cmpId)
	clearBit(&e.archetypeId, cmpId)
	t := db.ensureTable(e.archetypeId)
	e.idx = len(t.entries)
	t.entries = append(t.entries, e)
	return nil
}

func (db *database) deleteEntity(e *entity) {
	db.tables[e.archetypeId].entries = append(
		db.tables[e.archetypeId].entries[:e.idx],
		db.tables[e.archetypeId].entries[e.idx+1:]...,
	)
	delete(db.refs, e.id)
}

func (db *database) ensureTable(archetypeId archetypeId) *table {
	t, found := db.tables[archetypeId]
	if !found {
		t = &table{
			archetypeId: archetypeId,
			entries:     make([]*entity, 0),
		}
		db.tables[archetypeId] = t
	}
	return t
}

func (d *database) parseComponent(cmp any) (componentId, any) {
	cmpId := d.cm.ensureId(cmp)
	cmpValue := reflect.ValueOf(cmp)
	if cmpValue.Kind() != reflect.Pointer {
		ptr := reflect.New(cmpValue.Type())
		ptr.Elem().Set(cmpValue)
		cmp = ptr.Interface()
	}
	return cmpId, cmp
}
