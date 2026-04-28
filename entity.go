package main

import (
	"strings"
	"time"
)

// EntityCallback is a custom behavior function for one entity.
// It runs during animation updates and can move or change state.
type EntityCallback func(*Entity, *Animation) bool

// EntityCollisionHandler runs when an entity touches others.
// Use it to react to hits like fish meeting a hook or shark.
type EntityCollisionHandler func(*Entity, *Animation)

// EntityDeathHandler runs right before an entity is removed.
// It can spawn replacement entities or trigger follow-up effects.
type EntityDeathHandler func(*Entity, *Animation)

// Entity stores everything needed for one on-screen object.
// It includes sprite frames, movement data, collision data, and life rules.
// Fish, bubbles, hooks, and decorations all use this same base type.
type Entity struct {
	Name       string
	EntityType string

	Shapes []string
	Colors []string

	X float64
	Y float64
	Z float64

	Callback      EntityCallback
	CallbackArgs  any
	DieTime       *time.Time
	DieOffscreen  bool
	DieFrame      int
	DeathCallback EntityDeathHandler

	DefaultColor string
	Physical     bool
	CollHandler  EntityCollisionHandler
	AutoTrans    bool
	Transparent  rune

	CurrentFrame int
	FrameTime    float64
	FrameCount   int
	Collision    []*Entity
	Alive        bool
	width        int
	height       int
}

// NewEntityOptions is the input bundle used by NewEntity.
// Shape holds one or more animation frames; Color holds matching color masks.
// Single-frame entities pass a one-element slice; animated entities pass more.
// CallbackArgs can be []float64 movement data or a custom mode map.
type NewEntityOptions struct {
	Name          string
	EntityType    string
	Shape         []string
	Color         []string
	Position      [3]int
	Callback      EntityCallback
	CallbackArgs  any
	DieTime       *time.Time
	DieOffscreen  bool
	DieFrame      int
	DeathCallback EntityDeathHandler
	DefaultColor  string
	Physical      bool
	CollHandler   EntityCollisionHandler
	AutoTrans     bool
}

// NewEntity builds an Entity from options and fills safe defaults.
// It normalizes colors, shape slices, and initial callback arguments.
// This gives all entities a consistent starting state.
func NewEntity(opts NewEntityOptions) *Entity {
	e := &Entity{
		Name:          opts.Name,
		EntityType:    opts.EntityType,
		Shapes:        opts.Shape,
		Colors:        opts.Color,
		X:             float64(opts.Position[0]),
		Y:             float64(opts.Position[1]),
		Z:             float64(opts.Position[2]),
		Callback:      opts.Callback,
		CallbackArgs:  opts.CallbackArgs,
		DieTime:       opts.DieTime,
		DieOffscreen:  opts.DieOffscreen,
		DieFrame:      opts.DieFrame,
		DeathCallback: opts.DeathCallback,
		DefaultColor:  strings.ToUpper(opts.DefaultColor),
		Physical:      opts.Physical,
		CollHandler:   opts.CollHandler,
		AutoTrans:     opts.AutoTrans,
		Transparent:   ' ',
		Alive:         true,
	}
	if e.DefaultColor == "" {
		e.DefaultColor = "WHITE"
	}
	if len(e.Shapes) == 0 {
		e.Shapes = []string{""}
	}
	if len(e.Colors) == 0 {
		e.Colors = []string{""}
	}
	if e.CallbackArgs == nil {
		e.CallbackArgs = []float64{0, 0, 0, 0.5}
	}
	e.updateDimensions()
	return e
}


// updateDimensions recalculates current frame width and height.
// The renderer and collision checks use these values every frame.
// It should run whenever shape data may change.
func (e *Entity) updateDimensions() {
	lines := strings.Split(e.Shapes[0], "\n")
	e.height = len(lines)
	maxW := 0
	for _, ln := range lines {
		if len(ln) > maxW {
			maxW = len(ln)
		}
	}
	e.width = maxW
}

// Position returns integer coordinates for drawing and collisions.
// The entity stores float movement internally, then rounds by cast.
func (e *Entity) Position() (int, int, int) {
	return int(e.X), int(e.Y), int(e.Z)
}

// Size returns current sprite width and height in cells.
// This helps clipping and collision math stay simple.
func (e *Entity) Size() (int, int) {
	return e.width, e.height
}

// CurrentShape picks the frame to draw right now.
// It loops automatically when frame index passes frame count.
func (e *Entity) CurrentShape() string {
	if len(e.Shapes) == 0 {
		return ""
	}
	return e.Shapes[e.CurrentFrame%len(e.Shapes)]
}

// CurrentColor picks the active color-mask frame.
// Like shapes, this cycles through available mask frames.
func (e *Entity) CurrentColor() string {
	if len(e.Colors) == 0 {
		return ""
	}
	return e.Colors[e.CurrentFrame%len(e.Colors)]
}

// MoveEntity is the default movement logic when no custom callback exists.
// For []float64 args, order is [dx, dy, dz, frameStep].
// frameStep accumulates until >= 1, then frame index is advanced.
func (e *Entity) MoveEntity(_ *Animation) bool {
	switch args := e.CallbackArgs.(type) {
	case []float64:
		if len(args) >= 3 {
			e.X += args[0]
			e.Y += args[1]
			e.Z += args[2]
		}
		if len(args) >= 4 && args[3] > 0 {
			e.FrameTime += args[3]
			if e.FrameTime >= 1 {
				e.CurrentFrame++
				e.FrameTime = 0
			}
			e.FrameCount++
		}
	default:
		if len(e.Shapes) > 1 {
			e.FrameTime += 0.1
			if e.FrameTime >= 1 {
				e.CurrentFrame++
				e.FrameTime = 0
				e.FrameCount++
			}
		}
	}
	return true
}

// Kill marks the entity as dead for cleanup.
// The animation loop removes dead entities on next update.
func (e *Entity) Kill() {
	e.Alive = false
}

// ShouldDie checks all removal rules for this entity.
// It handles manual kill, time/frame limits, and offscreen cleanup.
// Returning true means the entity should be deleted now.
func (e *Entity) ShouldDie(screenWidth, screenHeight int, now time.Time) bool {
	if !e.Alive {
		return true
	}
	if e.DieTime != nil && now.After(*e.DieTime) {
		return true
	}
	if e.DieFrame > 0 && e.FrameCount >= e.DieFrame {
		return true
	}
	if e.DieOffscreen {
		if e.X+float64(e.width) < 0 || e.X >= float64(screenWidth) || e.Y+float64(e.height) < 0 || e.Y >= float64(screenHeight) {
			return true
		}
	}
	return false
}

// Update runs this entity's behavior for one frame.
// First it moves (custom callback or default movement), then handles hits.
// Collision handler runs only after collision lists were prepared by Animation.
func (e *Entity) Update(anim *Animation) {
	if e.Callback != nil {
		e.Callback(e, anim)
	} else {
		e.MoveEntity(anim)
	}
	if e.CollHandler != nil && len(e.Collision) > 0 {
		e.CollHandler(e, anim)
	}
}
