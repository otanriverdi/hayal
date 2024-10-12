package ecs

type SystemCtx interface {
	Spawn(cmp any) (entity, error)
	Destroy(entity entity) error
	AddComponent(entity entity, cmp any) error
	RemoveComponent(entity entity, cmp any) error
	Query(cmps ...any) (func(yield func(QueryResult) bool), error)
}
