package main

import (
	"math/rand"
	"strings"
)

// AddShark creates two linked entities: shark art and teeth hitbox.
// Keeping teeth separate makes collision checks simple and precise.
// Shark death callback later removes teeth and starts next event.
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

// SharkDeath removes teeth when shark exits.
// It then spawns another random special object.
func SharkDeath(_ *Entity, anim *Animation) {
	for _, t := range anim.GetEntitiesByType("teeth") {
		anim.DelEntity(t)
	}
	RandomObject(nil, anim)
}

// AddShip spawns a boat moving on the surface layer.
// It includes color mask so sails and hull have detail.
func AddShip(_ *Entity, anim *Animation) {
	shapes := [2]string{
		`     |    |    |
    )_)  )_)  )_)
   )___))___))___)\
  )____)____)_____)\\
_____|____|____|____\\\__
\                   /`,
		`         |    |    |
        (_(  (_(  (_(
      /(___((___((___(
    //(_____(____(____(
__///____|____|____|_____
    \                   /`,
	}
	colors := [2]string{
		`     y    y    y

                  w
                   ww
yyyyyyyyyyyyyyyyyyyywwwyy
y                   y`,
		`         y    y    y

      w
    ww
yywwwyyyyyyyyyyyyyyyyyyyy
    y                   y`,
	}
	dir := rand.Intn(2)
	speed := 1.0
	x := -24
	if dir == 1 {
		speed = -1
		x = anim.Width() - 2
	}
	anim.NewEntity(NewEntityOptions{
		Shape:         shapes[dir],
		Color:         colors[dir],
		AutoTrans:     true,
		Position:      [3]int{x, 0, Depth["water_gap1"]},
		DefaultColor:  "WHITE",
		CallbackArgs:  []float64{speed, 0, 0, 0},
		DieOffscreen:  true,
		DeathCallback: RandomObject,
	})
}

// AddWhale builds animation frames in code instead of hardcoding all frames.
// We start with idle body frames, then append aligned spout variations.
// Same color masks are reused while shape frames change over time.
func AddWhale(_ *Entity, anim *Animation) {
	shapes := [2]string{
		`        .-----:
      .'       ` + "`" + `.
,    /       (o) \
\` + "`" + `._/          ,__)`,
		`    :-----.
  .'       ` + "`" + `.
 / (o)       \    ,
(__,          \_.'/`,
	}
	colors := [2]string{
		`             C C
           CCCCCCC
           C  C  C
        BBBBBBB
      BB       BB
B    B       BWB B
BBBBB          BBBB`,
		`   C C
 CCCCCCC
 C  C  C
    BBBBBBB
  BB       BB
 B BWB       B    B
BBBB          BBBBB`,
	}
	waterSpouts := []string{
		"\n\n\n   :",
		"\n\n   :\n   :",
		"\n  . .\n  -:-\n   :",
		"\n  . .\n .-:-.\n   :",
		"\n  . .\n'.-:-.`\n'  :  '",
		"\n\n .- -.\n;  :  ;",
		"\n\n\n;     ;",
	}
	dir := rand.Intn(2)
	speed := 0.5
	x := -18
	spoutAlign := 11
	if dir == 1 {
		speed = -0.5
		x = anim.Width() - 2
		spoutAlign = 1
	}

	whaleAnim := make([]string, 0, 12)
	whaleAnimMask := make([]string, 0, 12)
	for i := 0; i < 5; i++ {
		whaleAnim = append(whaleAnim, "\n\n\n"+shapes[dir])
		whaleAnimMask = append(whaleAnimMask, colors[dir])
	}
	for _, spoutFrame := range waterSpouts {
		spoutLines := strings.Split(spoutFrame, "\n")
		sep := "\n" + strings.Repeat(" ", spoutAlign)
		alignedSpout := strings.Join(spoutLines, sep)
		whaleAnim = append(whaleAnim, alignedSpout+"\n"+shapes[dir])
		whaleAnimMask = append(whaleAnimMask, colors[dir])
	}

	anim.NewEntity(NewEntityOptions{
		Shape:         whaleAnim,
		Color:         whaleAnimMask,
		AutoTrans:     true,
		Position:      [3]int{x, 0, Depth["water_gap2"]},
		DefaultColor:  "WHITE",
		CallbackArgs:  []float64{speed, 0, 0, 1},
		DieOffscreen:  true,
		DeathCallback: RandomObject,
	})
}

// AddMonster spawns a sea monster crossing the screen.
// It is a large special event creature in deeper layer.
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

// AddBigFish creates a larger fast fish variant.
// It appears as a random special object.
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

// AddFishhook creates a 3-part system: line, visible hook, and catch point.
// All parts share the same callback mode so they move together.
// Hook point is the physical part that fish collides with.
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

// FishhookCallback acts like a tiny state machine.
// "lowering" moves down to max depth; "hooked" reels upward to top clamp.
// This callback uses mode map args, unlike most entities' []float64 args.
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

// Retract switches an entity into "hooked" upward movement.
// Used for fish, line, and hook after a catch event.
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

// GroupDeath removes all entities of listed types from scene.
// It is used for grouped cleanup (for example hook + line + point).
// After cleanup, it chains into the next random event.
func GroupDeath(entity *Entity, anim *Animation, boundTypes []string) {
	for _, tp := range boundTypes {
		for _, obj := range anim.GetEntitiesByType(tp) {
			anim.DelEntity(obj)
		}
	}
	RandomObject(entity, anim)
}

// AddDucks spawns animated ducks on the water surface.
// They cycle wing/pose frames while moving sideways.
func AddDucks(_ *Entity, anim *Animation) {
	dir := rand.Intn(2)
	shapes := [2][]string{
		{
			`      _          _          _
,____(')=  ,____(')=  ,____(')<
 \~~= ')    \~~= ')    \~~= ')`,
			`      _          _          _
,____(')=  ,____(')<  ,____(')=
 \~~= ')    \~~= ')    \~~= ')`,
			`      _          _          _
,____(')<  ,____(')=  ,____(')=
 \~~= ')    \~~= ')    \~~= ')`,
		},
		{
			`  _          _          _
>(')____,  =(')____,  =(')____,
 (` + "`" + ` =~~/    (` + "`" + ` =~~/    (` + "`" + ` =~~/`,
			`  _          _          _
=(')____,  >(')____,  =(')____,
 (` + "`" + ` =~~/    (` + "`" + ` =~~/    (` + "`" + ` =~~/`,
			`  _          _          _
=(')____,  =(')____,  >(')____,
 (` + "`" + ` =~~/    (` + "`" + ` =~~/    (` + "`" + ` =~~/`,
		},
	}
	colors := [2]string{
		`      g          g          g
wwwwwgcgy  wwwwwgcgy  wwwwwgcgy
 wwww Ww    wwww Ww    wwww Ww`,
		`  g          g          g
ygcgwwwww  ygcgwwwww  ygcgwwwww
 wW wwww    wW wwww    wW wwww`,
	}
	speed := 1.0
	x := -30
	if dir == 1 {
		speed = -1
		x = anim.Width() - 2
	}
	anim.NewEntity(NewEntityOptions{
		Shape:         shapes[dir],
		Color:         colors[dir],
		AutoTrans:     true,
		Position:      [3]int{x, 5, Depth["water_gap3"]},
		CallbackArgs:  []float64{speed, 0, 0, 0.25},
		DeathCallback: RandomObject,
		DieOffscreen:  true,
		DefaultColor:  "WHITE",
	})
}

// AddDolphins spawns three dolphins with fixed spacing.
// Only the lead dolphin has death callback to avoid triple respawns.
// Followers are visual companions in the same formation.
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
	shapes := [2][]string{
		{
			"        ,\n      __)\\\n(\\_.-'    a`-.\n(/~~````(/~^^`",
			"        ,\n(\\__  __)\\\n(/~.''    a`-.\n    ````\\)~^^`",
		},
		{
			"     ,\n   _/(__\n.-'a    `-._/)\n'^^~\\)''''~~\\)",
			"     ,\n   _/(__  __/)\n.-'a    ``.~\\)\n'^^~(/''''",
		},
	}
	colors := [2]string{
		"\n\n\n          W",
		"\n\n\n   W",
	}
	for i := 0; i < 3; i++ {
		deathCb := EntityDeathHandler(nil)
		if i == 0 {
			deathCb = func(_ *Entity, a *Animation) { RandomObject(nil, a) }
		}
		anim.NewEntity(NewEntityOptions{
			Shape:     shapes[dir],
			Color:     colors[dir],
			AutoTrans: true,
			Position:  [3]int{x - (distance * (2 - i)), 5, Depth["water_gap3"]},
			CallbackArgs: []float64{
				speed, 0, 0, 0.5,
			},
			DeathCallback: deathCb,
			DieOffscreen:  true,
			DefaultColor:  "CYAN",
		})
	}
}

// AddSwan spawns a swan gliding near the top water line.
// Direction and sprite frame are picked randomly.
func AddSwan(_ *Entity, anim *Animation) {
	shapes := [2]string{
		`       ___
,_    / _,\
| \   \( \|
|  \_  \\
(_   \_) \
(\_   ` + "`" + `   \
 \   -=~  /`,
		` ___
/,_ \    _,
|/ )/   / |
  //  _/  |
 / ( /   _)
/   ` + "`" + `   _/)
\  ~=-   /`,
	}
	colors := [2]string{
		"\n\n         g\n         yy\n\n\n\n",
		"\n\n g\nyy\n\n\n\n",
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
		Color:         colors[dir],
		AutoTrans:     true,
		Position:      [3]int{x, 1, Depth["water_gap3"]},
		CallbackArgs:  []float64{speed, 0, 0, 0.25},
		DeathCallback: RandomObject,
		DieOffscreen:  true,
		DefaultColor:  "WHITE",
	})
}

// RandomObject is the special-event router for this game.
// Many death callbacks call this, so events form a continuous chain.
// Random choice keeps the aquarium from repeating one pattern.
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

// maxInt returns the larger of two integers.
// It is a small helper for safe random ranges.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
