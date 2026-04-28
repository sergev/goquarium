package main

import (
	"strings"
	"time"
)

type EntityCallback func(*Entity, *Animation) bool
type EntityCollisionHandler func(*Entity, *Animation)
type EntityDeathHandler func(*Entity, *Animation)

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

type NewEntityOptions struct {
	Name          string
	EntityType    string
	Shape         any
	Color         any
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

func NewEntity(opts NewEntityOptions) *Entity {
	e := &Entity{
		Name:          opts.Name,
		EntityType:    opts.EntityType,
		Shapes:        asFrameSlice(opts.Shape),
		Colors:        asFrameSlice(opts.Color),
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

func asFrameSlice(v any) []string {
	switch t := v.(type) {
	case nil:
		return nil
	case string:
		return []string{strings.ReplaceAll(t, "?", " ")}
	case []string:
		out := make([]string, 0, len(t))
		for _, s := range t {
			out = append(out, strings.ReplaceAll(s, "?", " "))
		}
		return out
	default:
		return []string{""}
	}
}

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

func (e *Entity) Position() (int, int, int) {
	return int(e.X), int(e.Y), int(e.Z)
}

func (e *Entity) Size() (int, int) {
	return e.width, e.height
}

func (e *Entity) CurrentShape() string {
	if len(e.Shapes) == 0 {
		return ""
	}
	return e.Shapes[e.CurrentFrame%len(e.Shapes)]
}

func (e *Entity) CurrentColor() string {
	if len(e.Colors) == 0 {
		return ""
	}
	return e.Colors[e.CurrentFrame%len(e.Colors)]
}

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

func (e *Entity) Kill() {
	e.Alive = false
}

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
