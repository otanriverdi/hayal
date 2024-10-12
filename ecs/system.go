package ecs

type SystemCtx interface {
	// Spawn initializes a new entity with the passed in component.
	Spawn(cmp any) (entity, error)
	// Destroy de-initializes the passed in entity.
	Destroy(entity entity) error
	// AddComponent adds the passed in component to the entity.
	AddComponent(entity entity, cmp any) error
	// RemoveComponent removes the passed in component from the entity.
	RemoveComponent(entity entity, cmp any) error
	// Query returns an iterator of all entities that includes the passed in components.
	Query(cmps ...any) (func(yield func(QueryResult) bool), error)
}
