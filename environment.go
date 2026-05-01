package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// AddEnvironment draws repeated water bands near the surface.
// They are also marked physical so bubbles can collide and pop there.
// This gives both visuals and simple "surface" collision behavior.
func AddEnvironment(anim *Animation) {
	waterSegments := []string{
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		"^^^^ ^^^  ^^^   ^^^    ^^^^      ",
		"^^^^      ^^^^     ^^^    ^^     ",
		"^^      ^^^^      ^^^    ^^^^^^  ",
	}
	segSize := len(waterSegments[0])
	repeat := (anim.Width() / segSize) + 1
	for i, seg := range waterSegments {
		tiled := strings.Repeat(seg, repeat)
		depthKey := fmt.Sprintf("water_line%d", i)
		depth, ok := Depth[depthKey]
		if !ok {
			depth = Depth["water_line0"]
		}
		anim.NewEntity(NewEntityOptions{
			Name:         fmt.Sprintf("water_seg_%d", i),
			EntityType:   "waterline",
			Shape:        []string{tiled},
			Position:     [3]int{0, i + 5, depth},
			DefaultColor: "CYAN",
			Physical:     true,
		})
	}
}

// AddCastle places a decorative castle near the bottom-right.
// It is static scenery and does not move over time.
func AddCastle(anim *Animation) {
	castleShape := `               T~~
               |
              /^\
             /   \
 _   _   _  /     \  _   _   _
[ ]_[ ]_[ ]/ _   _ \[ ]_[ ]_[ ]
|_=__-_ =_|_[ ]_[ ]_|_=-___-__|
 | _- =  | =_ = _    |= _=   |
 |= -[]  |- = _ =    |_-=_[] |
 | =_    |= - ___    | =_ =  |
 |=  []- |-  /| |\   |=_ =[] |
 |- =_   | =| | | |  |- = -  |
 |_______|__|_|_|_|__|_______|`

	castleColor := `                RR
                W
              Wyyw
             y   y
 W   W   W  yWWWWWy  W   W   W
WW WW WW WW W   W WwWW WW WW WW
WWWWWWW WWWWW W W WWWWWWWWWWWWWW
 W W W  W W W W W    W  W   WWW
 W  W   W  W W W     W W W  WWW
 W  W   W  W WWW     W W W  WWW
 W  W   W  W W W W   W  W   WWW
 W  W   W W W W W W  W  W   WWW
 WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW`
	anim.NewEntity(NewEntityOptions{
		EntityType:   "castle",
		Name:         "castle",
		Shape:        []string{castleShape},
		Color:        []string{castleColor},
		Position:     [3]int{anim.Width() - 32, anim.Height() - 13, Depth["castle"]},
		DefaultColor: "DARK_GREY",
	})
}

// AddSeaweed creates one seaweed with two alternating frames.
// It does not move position; only frame index changes to fake swaying.
// When lifetime ends, death callback spawns a replacement plant.
func AddSeaweed(_ *Entity, anim *Animation) {
	frames := []string{"", ""}
	height := rand.Intn(4) + 3
	for i := 1; i <= height; i++ {
		leftSide := i % 2
		rightSide := 1 - leftSide
		frames[leftSide] += "(\n"
		frames[rightSide] += " )\n"
	}
	maxX := anim.Width() - 2
	if maxX < 1 {
		maxX = 1
	}
	x := rand.Intn(maxX) + 1
	y := anim.Height() - height
	if y < 9 {
		y = 9
	}
	speed := 0.25 + (rand.Float64() * 0.05)
	life := time.Now().Add(time.Duration(8*60+rand.Intn(4*60)) * time.Second)
	anim.NewEntity(NewEntityOptions{
		Name:          fmt.Sprintf("seaweed_%f", rand.Float64()),
		EntityType:    "seaweed",
		Shape:         frames,
		Position:      [3]int{x, y, Depth["seaweed"]},
		CallbackArgs:  []float64{0, 0, 0, speed},
		DieTime:       &life,
		DeathCallback: AddSeaweed,
		DefaultColor:  "GREEN",
	})
}

// AddAllSeaweed seeds the initial seaweed population.
// Count uses a simple width/15 density rule for balanced scenery.
func AddAllSeaweed(anim *Animation) {
	count := anim.Width() / 15
	for i := 0; i < count; i++ {
		AddSeaweed(nil, anim)
	}
}
