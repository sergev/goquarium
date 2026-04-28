package main

import (
	"math/rand"
)

func AddShark(_ *Entity, anim *Animation) {
	shapes := [2]string{
		`   __
__/o )>`,
		`< ( o\__
    __`,
	}
	direction := rand.Intn(2)
	speed := 2.0
	x := -8
	if direction == 1 {
		speed = -2.0
		x = anim.Width() - 2
	}
	y := 9
	if anim.Height() > 19 {
		y = 9 + rand.Intn(anim.Height()-18)
	}
	teethX := x + 2
	if direction == 1 {
		teethX = x - 1
	}
	anim.NewEntity(NewEntityOptions{
		EntityType:   "teeth",
		Shape:        "*",
		Position:     [3]int{teethX, y + 1, Depth["shark"] + 1},
		CallbackArgs: []float64{speed, 0, 0},
		Physical:     true,
	})
	anim.NewEntity(NewEntityOptions{
		EntityType:    "shark",
		Shape:         shapes[direction],
		AutoTrans:     true,
		Position:      [3]int{x, y, Depth["shark"]},
		DefaultColor:  "CYAN",
		CallbackArgs:  []float64{speed, 0, 0},
		DieOffscreen:  true,
		DeathCallback: SharkDeath,
	})
}

func SharkDeath(_ *Entity, anim *Animation) {
	for _, t := range anim.GetEntitiesByType("teeth") {
		anim.DelEntity(t)
	}
	RandomObject(nil, anim)
}

func AddShip(_ *Entity, anim *Animation) {
	shapes := [2]string{
		`   |\n __|___\n \_____/`,
		` |\n___|__\n\_____/`,
	}
	dir := rand.Intn(2)
	speed := 1.0
	x := -10
	if dir == 1 {
		speed = -1
		x = anim.Width() - 2
	}
	anim.NewEntity(NewEntityOptions{
		Shape:         shapes[dir],
		AutoTrans:     true,
		Position:      [3]int{x, 0, Depth["water_gap1"]},
		DefaultColor:  "WHITE",
		CallbackArgs:  []float64{speed, 0, 0, 0},
		DieOffscreen:  true,
		DeathCallback: RandomObject,
	})
}

func AddWhale(_ *Entity, anim *Animation) {
	shapes := [2]string{
		`  .-----:
.' (o)   \`,
		`:-----.
/   (o) '.`,
	}
	dir := rand.Intn(2)
	speed := 0.5
	x := -12
	if dir == 1 {
		speed = -0.5
		x = anim.Width() - 2
	}
	anim.NewEntity(NewEntityOptions{
		Shape:         shapes[dir],
		AutoTrans:     true,
		Position:      [3]int{x, 0, Depth["water_gap2"]},
		DefaultColor:  "WHITE",
		CallbackArgs:  []float64{speed, 0, 0, 1},
		DieOffscreen:  true,
		DeathCallback: RandomObject,
	})
}

func AddMonster(_ *Entity, anim *Animation) {
	shapes := [2]string{
		` _a_a
/ oo\____`,
		`____/oo \
 a_a_`,
	}
	dir := rand.Intn(2)
	speed := 2.0
	x := -10
	if dir == 1 {
		speed = -2
		x = anim.Width() - 2
	}
	anim.NewEntity(NewEntityOptions{
		Shape:         shapes[dir],
		AutoTrans:     true,
		Position:      [3]int{x, 2, Depth["water_gap2"]},
		CallbackArgs:  []float64{speed, 0, 0, 0.25},
		DeathCallback: RandomObject,
		DieOffscreen:  true,
		DefaultColor:  "GREEN",
	})
}

func AddBigFish(_ *Entity, anim *Animation) {
	shapes := [2]string{
		` __
<` + "`" + `)))><`,
		`><(((´>
  __`,
	}
	dir := rand.Intn(2)
	speed := 2.5
	x := -8
	if dir == 1 {
		speed = -2.5
		x = anim.Width() - 1
	}
	y := 9
	if anim.Height()-14 > 9 {
		y = 9 + rand.Intn(anim.Height()-14-9+1)
	}
	anim.NewEntity(NewEntityOptions{
		Shape:         shapes[dir],
		AutoTrans:     true,
		Position:      [3]int{x, y, Depth["shark"]},
		CallbackArgs:  []float64{speed, 0, 0},
		DeathCallback: RandomObject,
		DieOffscreen:  true,
		DefaultColor:  "YELLOW",
	})
}

func AddFishhook(_ *Entity, anim *Animation) {
	x := 10 + rand.Intn(maxInt(1, anim.Width()-30))
	yStart := -20
	anim.NewEntity(NewEntityOptions{
		EntityType:   "fishline",
		Shape:        "|\n|\n|\n|\n|\n|\n|\n|\n|\n|\n",
		Position:     [3]int{x + 7, yStart - 10, Depth["water_line1"]},
		AutoTrans:    true,
		Callback:     FishhookCallback,
		CallbackArgs: map[string]string{"mode": "lowering"},
	})
	anim.NewEntity(NewEntityOptions{
		EntityType:    "fishhook",
		Shape:         "o\n||\n\\\\//",
		Position:      [3]int{x, yStart, Depth["water_line1"]},
		AutoTrans:     true,
		DieOffscreen:  true,
		DefaultColor:  "GREEN",
		Callback:      FishhookCallback,
		CallbackArgs:  map[string]string{"mode": "lowering"},
		DeathCallback: func(ent *Entity, a *Animation) { GroupDeath(ent, a, []string{"hook_point", "fishline"}) },
	})
	anim.NewEntity(NewEntityOptions{
		EntityType:   "hook_point",
		Shape:        ".\n \n\\\n ",
		Position:     [3]int{x + 1, yStart + 2, Depth["shark"] + 1},
		Physical:     true,
		DefaultColor: "GREEN",
		Callback:     FishhookCallback,
		CallbackArgs: map[string]string{"mode": "lowering"},
	})
}

func FishhookCallback(entity *Entity, anim *Animation) bool {
	mode := ""
	switch args := entity.CallbackArgs.(type) {
	case map[string]string:
		mode = args["mode"]
	}
	if mode == "hooked" {
		entity.Y -= 2
		if entity.Y < -10 {
			entity.Y = -10
		}
		return true
	}
	maxDepth := int(float64(anim.Height()) * 0.75)
	if int(entity.Y) < maxDepth {
		entity.Y += 2
	} else {
		entity.Y = float64(maxDepth)
	}
	return true
}

func Retract(entity *Entity, _ *Animation) {
	entity.Physical = false
	if entity.EntityType == "fish" {
		entity.Z = float64(Depth["water_gap2"])
		entity.Callback = FishhookCallback
		entity.CallbackArgs = map[string]string{"mode": "hooked"}
		return
	}
	entity.CallbackArgs = map[string]string{"mode": "hooked"}
}

func GroupDeath(entity *Entity, anim *Animation, boundTypes []string) {
	for _, tp := range boundTypes {
		for _, obj := range anim.GetEntitiesByType(tp) {
			anim.DelEntity(obj)
		}
	}
	RandomObject(entity, anim)
}

func AddDucks(_ *Entity, anim *Animation) {
	dir := rand.Intn(2)
	shape := []string{
		`,____(')=  ,____(')=  ,____(')<`,
		`>(')____,  =(')____,  =(')____,`,
	}[dir]
	speed := 1.0
	x := -30
	if dir == 1 {
		speed = -1
		x = anim.Width() - 2
	}
	anim.NewEntity(NewEntityOptions{
		Shape:         shape,
		AutoTrans:     true,
		Position:      [3]int{x, 5, Depth["water_gap3"]},
		CallbackArgs:  []float64{speed, 0, 0, 0.25},
		DeathCallback: RandomObject,
		DieOffscreen:  true,
		DefaultColor:  "WHITE",
	})
}

func AddDolphins(_ *Entity, anim *Animation) {
	dir := rand.Intn(2)
	speed := 2.0
	x := -13
	distance := 15
	if dir == 1 {
		speed = -2
		x = anim.Width() - 2
		distance = -15
	}
	shape := []string{
		"  __)\n(/_.-'`",
		" _/(__\n.-'a  `-._/)",
	}[dir]
	for i := 0; i < 3; i++ {
		anim.NewEntity(NewEntityOptions{
			Shape:     shape,
			AutoTrans: true,
			Position:  [3]int{x - (distance * (2 - i)), 5, Depth["water_gap3"]},
			CallbackArgs: []float64{
				speed, 0, 0, 0.5,
			},
			DeathCallback: func(_ *Entity, a *Animation) {
				if i == 0 {
					RandomObject(nil, a)
				}
			},
			DieOffscreen: true,
			DefaultColor: "CYAN",
		})
	}
}

func AddSwan(_ *Entity, anim *Animation) {
	shapes := [2]string{
		`  ___
,_/ _ \`,
		` ___
/ _ \_,`,
	}
	dir := rand.Intn(2)
	speed := 1.0
	x := -10
	if dir == 1 {
		speed = -1
		x = anim.Width() - 2
	}
	anim.NewEntity(NewEntityOptions{
		Shape:         shapes[dir],
		AutoTrans:     true,
		Position:      [3]int{x, 1, Depth["water_gap3"]},
		CallbackArgs:  []float64{speed, 0, 0, 0.25},
		DeathCallback: RandomObject,
		DieOffscreen:  true,
		DefaultColor:  "WHITE",
	})
}

func RandomObject(dead *Entity, anim *Animation) {
	randomObjects := []func(*Entity, *Animation){
		AddShip,
		AddWhale,
		AddMonster,
		AddBigFish,
		AddShark,
		AddFishhook,
		AddSwan,
		AddDucks,
		AddDolphins,
	}
	spawner := randomObjects[rand.Intn(len(randomObjects))]
	spawner(dead, anim)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
