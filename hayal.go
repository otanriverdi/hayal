// Package hayal provides an experimental ecs based game engine.
//
//
// To the Game struct, you can assign systems to be executed in various steps of the game engine, that would initialize
// it's own entities and update logic:
//
//	game := hayal.New()
//	game.AddSystem(hayal.GameLoopStepUpdate, system)
//	game.Run()
//
// The system functions have access to a gameCtx that can be used to interface with the game world.
//
//  func System(ctx *gameCtx) {
//    e, err := ctx.Spawn(transform{x: 5, y: 10})
//    if err != nil {
//      return err
//    }
//
//    err = ctx.Destroy(e)
//    if err != nil {
//      return err
//    }
//
//    err = ctx.AddComponent(e, 5)
//    if err != nil {
//      return err
//    }
//
//    err = ctx.RemoveComponent(e, 5)
//    if err != nil {
//      return err
//    }
//
//    iter, err := ecs.Query(transform{})
//    if err != nil {
//      return err
//    }
//    for entity := range iter {
//      cmp, err := GetComponent[transform](&entity)
//      if err != nil {
//        return err
//      }
//
//      err = SetComponent(&entity, transform{x: 20, y: 23})
//      if err != nil {
//        return err
//      }
//    }
//
//    ctx.Exit()
//  }
//
package hayal

import (
	"sync"

	"github.com/otanriverdi/hayal/ecs"
)

type gameLoopStep = uint

const (
	// GameLoopStepInit runs once on game start.
	GameLoopStepInit gameLoopStep = iota
	// GameLoopStepPreUpdate runs first on every tick. This is where we update resources that would be used in the
	// update stage.
	GameLoopStepPreUpdate
	// GameLoopStateUpdate runs on every tick. This is where we update game play systems.
	GameLoopStateUpdate
	// GameLoopStateDraw runs on every tick after Update. This is where we render and apply component updates.
	GameLoopStateDraw
	// GameLoopStateDeinit runs once before exit. Use this for cleanup.
	GameLoopStateDeinit
)

type gameCtx struct {
	ecs.ECS
	exit chan struct{}
}

func (ctx *gameCtx) Exit() {
	close(ctx.exit)
}

type SystemCtx interface {
	ecs.SystemCtx
	// Exit signals the game to run cleanup and exit gracefully.
	Exit()
}

type System = func(ctx SystemCtx) error

type Game struct {
	ctx *gameCtx
	// len of first dimension matches step count
	schedule [5][]System
}

// New initializes a new game.
func New() Game {
	return Game{ctx: &gameCtx{ECS: ecs.New()}}
}

// AddSystem adds a system to the schedule to be executed in various steps of the game loop. Check GameLoopStep*
// constants for various steps and when they run.
func (g *Game) AddSystem(step gameLoopStep, system System) {
	if g.schedule[step] == nil {
		g.schedule[step] = make([]System, 0)
	}
	g.schedule[step] = append(g.schedule[step], system)
}

// Run starts the schedule and the execution of the game.
func (g *Game) Run() {
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

type Plugin = func(g *Game)

func (g *Game) Plug(plugin Plugin) {
	plugin(g)
}

func (g *Game) executeStep(step gameLoopStep) {
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

// GetComponent extracts a copy of component data from the passed in query result.
func GetComponent[C any](qr *ecs.QueryResult) (C, error) {
	return ecs.GetComponent[C](qr)
}

// SetComponent sets the component data to be the passed in component.
func SetComponent(qr *ecs.QueryResult, cmp any) error {
	return ecs.SetComponent(qr, cmp)
}
