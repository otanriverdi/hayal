package hayal

import (
	"sync"

	"github.com/otanriverdi/hayal/ecs"
)

type gameCtx struct {
	ecs.ECS
	exit chan struct{}
}

func (ctx *gameCtx) Exit() {
	close(ctx.exit)
}

type GameLoopStep = uint

const (
	GameLoopStepInit GameLoopStep = iota
	GameLoopStepPreUpdate
	GameLoopStateUpdate
	GameLoopStateDraw
	GameLoopStateDeinit
)

type systemCtx interface {
	ecs.SystemCtx
	Exit()
}

type System = func(ctx systemCtx) error

type game struct {
	ctx *gameCtx
	// len of first dimension matches step count
	schedule [5][]System
}

func New() game {
	e := ecs.New()
	return game{ctx: &gameCtx{ECS: e}}
}

func (g *game) AddSystem(step GameLoopStep, system System) {
	if g.schedule[step] == nil {
		g.schedule[step] = make([]System, 0)
	}
	g.schedule[step] = append(g.schedule[step], system)
}

func (g *game) Run() {
	g.executeStep(GameLoopStepInit)
	for {
		select {
		case <-g.ctx.exit:
			g.executeStep(GameLoopStateDeinit)
			return
		default:
			g.executeStep(GameLoopStepPreUpdate)
			g.executeStep(GameLoopStateUpdate)
			g.executeStep(GameLoopStateDraw)
		}
	}
}

func (g *game) executeStep(step GameLoopStep) {
	var wg sync.WaitGroup
	for _, sys := range g.schedule[GameLoopStepInit] {
		wg.Add(1)
		go func(sys System) {
			defer wg.Done()
			err := sys(g.ctx)
			if err != nil {
				panic(err)
			}
		}(sys)
	}
	wg.Wait()
}

func GetComponent[C any](qr *ecs.QueryResult) (C, error) {
	return ecs.GetComponent[C](qr)
}

func SetComponent(qr *ecs.QueryResult, cmp any) error {
	return ecs.SetComponent(qr, cmp)
}
