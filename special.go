package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func debugf(format string, args ...any) {
	if DebugLogEnabled {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

// AddShark creates two linked entities: shark art and teeth hitbox.
// Keeping teeth separate makes collision checks simple and precise.
// Shark death callback later removes teeth and starts next event.
func AddShark(_ *Entity, anim *Animation) {
	shapes := [2]string{
		"                              __\n                             ( `\\\n  ,                          )   `\\\n;' `.                       (     `\\__\n ;   `.             __..---''          `~~~~-._\n  `.   `.____...--''                       (b  `--._\n    >                     _.-'      .((      ._     )\n  .`.-`--...__         .-'     -.___.....-(|/|/|/|/'\n ;.'         `. ...----`.___.',,,_______......---'\n '           '-'",
		"                     __\n                    /' )\n                  /'   (                          ,\n              __/'     )                       .' `;\n      _.-~~~~'          ``---..__             .'   ;\n _.--'  b)                       ``--...____.'   .'\n(     _.      )).      `-._                     <\n `\\|\\|\\|\\|)-.....___.-     `-.         __...--'-.'.\n   `---......_______,,,`.___.'----... .'         `.;\n                                     `-`           `",
	}
	colors := [2]string{
		"\n\n\n\n\n                                           cR\n \n                                          cWWWWWWWW\n\n\n",
		"\n\n\n\n        Rc\n\n  WWWWWWWWc\n\n\n\n",
	}
	direction := rand.Intn(2)
	x := -53
	y := 9
	teethX := -9
	teethY := y + 7
	speed := 2.0
	if anim.Height() > 19 {
		// Mirror upstream placement to keep sharks in deeper water.
		y = 9 + rand.Intn(maxInt(1, anim.Height()-19)+1)
		teethY = y + 7
	}
	if direction == 1 {
		speed = -2.0
		x = anim.Width() - 2
		teethX = x + 9
	}
	anim.NewEntity(NewEntityOptions{
		EntityType:   "teeth",
		Shape:        []string{"*"},
		Position:     [3]int{teethX, teethY, Depth["shark"] + 1},
		CallbackArgs: []float64{speed, 0, 0},
		Physical:     true,
	})
	anim.NewEntity(NewEntityOptions{
		EntityType:    "shark",
		Shape:         []string{shapes[direction]},
		Color:         []string{colors[direction]},
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
		Shape:         []string{shapes[dir]},
		Color:         []string{colors[dir]},
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
// It picks between the newer and classic animated monster designs.
func AddMonster(old *Entity, anim *Animation) {
	if rand.Intn(2) == 0 {
		addNewMonster(old, anim)
		return
	}
	addOldMonster(old, anim)
}

// addNewMonster creates the larger modern monster variant.
// It uses two animation frames and keeps the eye highlight mask.
func addNewMonster(_ *Entity, anim *Animation) {
	shapes := [2][]string{
		{
			"\n         _   _                   _   _       _a_a\n       _{.`=`.}_     _   _     _{.`=`.}_    {/ ''\\_\n _    {.'  _  '.}   {.`'`.}   {.'  _  '.}  {|  ._oo)\n{ \\  {/  .'~'.  \\}  {/ .-. \\}  {/  .'~'.  \\} {/  |",
			"\n                      _   _                    _a_a\n  _      _   _     _{.`=`.}_     _   _      {/ ''\\_\n { \\    {.`'`.}   {.'  _  '.}   {.`'`.}    {|  ._oo)\n  \\ \\  {/ .-. \\}  {/  .'~'.  \\}  {/ .-. \\}   {/  |",
		},
		{
			"\n   a_a_       _   _                   _   _\n _/'' \\}    _{.`=`.}_     _   _     _{.`=`.}_\n(oo_.  |}  {.'  _  '.}   {.`'`.}   {.'  _  '.}    _\n    |  \\} {/  .'~'.  \\}  {/ .-. \\}  {/  .'~'.  \\}  / }",
			"\n   a_a_                    _   _\n _/'' \\}      _   _     _{.`=`.}_     _   _      _\n(oo_.  |}    {.`'`.}   {.'  _  '.}   {.`'`.}    / }\n    |  \\}   {/ .-. \\}  {/  .'~'.  \\}  {/ .-. \\}  / /",
		},
	}
	colors := [2]string{
		"\n                                                W W\n\n\n\n",
		"\n   W W\n\n\n\n",
	}
	dir := rand.Intn(2)
	speed := 2.0
	x := -54
	if dir == 1 {
		speed = -2.0
		x = anim.Width() - 2
	}
	anim.NewEntity(NewEntityOptions{
		Shape:         shapes[dir],
		Color:         []string{colors[dir], colors[dir]},
		AutoTrans:     true,
		Position:      [3]int{x, 2, Depth["water_gap2"]},
		CallbackArgs:  []float64{speed, 0, 0, 0.25},
		DeathCallback: RandomObject,
		DieOffscreen:  true,
		DefaultColor:  "GREEN",
	})
}

// addOldMonster creates the classic sea monster variant.
// This one uses four animation frames for the body wake motion.
func addOldMonster(_ *Entity, anim *Animation) {
	shapes := [2][]string{
		{
			"\n                                                          ____\n            __                                          /   o  \\\n          /    \\        _                     _       /     ____ >\n  _      |  __  |     /   \\        _        /   \\   |     |\n | \\     |  ||  |    |     |     /   \\    |     |  |     |",
			"\n                                                          ____\n                                             __         /   o  \\\n             _                     _       /    \\     /     ____ >\n   _       /   \\        _        /   \\   |  __  |   |     |\n  | \\     |     |     /   \\    |     |  |  ||  |   |     |",
			"\n                                                          ____\n                                  __                  /   o  \\\n _                      _       /    \\        _     /     ____ >\n| \\          _        /   \\   |  __  |     /   \\  |     |\n \\ \\       /   \\    |     |  |  ||  |    |     | |     |",
			"\n                                                          ____\n                       __                             /   o  \\\n  _          _       /    \\        _                /     ____ >\n | \\       /   \\   |  __  |     /   \\        _    |     |\n  \\ \\     |     |  |  ||  |    |     |     /   \\  |     |",
		},
		{
			"\n    ____\n  /  o   \\                                          __\n< ____     \\       _                     _        /    \\\n      |     |   /   \\        _        /   \\     |  __  |      _\n      |     |  |     |     /   \\    |     |    |  ||  |     / |",
			"\n    ____\n  /  o   \\         __\n< ____     \\     /    \\       _                     _\n      |     |   |  __  |    /   \\        _        /   \\       _\n      |     |   |  ||  |   |     |     /   \\     |     |     / |",
			"\n    ____\n  /  o   \\                  __\n< ____     \\     _        /    \\       _                      _\n      |     |  /   \\     |  __  |   /   \\        _          / |\n      |     | |     |    |  ||  |  |     |    /   \\       / /",
			"\n    ____\n  /  o   \\                             __\n< ____     \\                _        /    \\       _          _\n      |     |    _        /   \\     |  __  |   /   \\       / |\n      |     |  /   \\    |     |    |  ||  |  |     |     / /",
		},
	}
	colors := [2]string{
		"\n\n                                                            W\n\n\n",
		"\n\n     W\n\n\n",
	}
	dir := rand.Intn(2)
	speed := 2.0
	x := -64
	if dir == 1 {
		speed = -2.0
		x = anim.Width() - 2
	}
	anim.NewEntity(NewEntityOptions{
		Shape:         shapes[dir],
		Color:         []string{colors[dir], colors[dir], colors[dir], colors[dir]},
		AutoTrans:     true,
		Position:      [3]int{x, 2, Depth["water_gap2"]},
		CallbackArgs:  []float64{speed, 0, 0, 0.25},
		DeathCallback: RandomObject,
		DieOffscreen:  true,
		DefaultColor:  "GREEN",
	})
}

// AddBigFish creates a larger fast fish variant.
// It keeps upstream weighting: design2 appears 2/3 of the time.
func AddBigFish(_ *Entity, anim *Animation) {
	if rand.Intn(3) > 0 {
		addBigFish2(nil, anim)
		return
	}
	addBigFish1(nil, anim)
}

func addBigFish1(_ *Entity, anim *Animation) {
	shapes := [2]string{
		" ______\n`\"\".  `````-----.....__\n     `.  .      .       `-.\n       :     .     .       `.\n ,     :   .    .          _ :\n: `.   :                  (@) `._\n `. `..'     .     =`-.       .__)\n   ;     .        =  ~  :     .-\"\n .' .'`.   .    .  =.-'  `._ .'\n: .'   :               .   .'\n '   .'  .    .     .   .-'\n   .'____....----''.'=.'\n   \"\"             .'.'\n               ''\"'`",
		"                           ______\n          __.....-----'''''  .-\"\"'\n       .-'       .      .  .'\n     .'       .     .     :\n    : _          .    .   :     ,\n _.' (@)                  :   .' :\n(__.       .-'=     .     `..' .'\n \"-.     :  ~  =        .     ;\n   `. _.'  `-.=  .    .   .'`. `.\n     `.   .               :   `. :\n       `-.   .     .    .  `.   `\n          `.=`.``----....____`.\n            `.`.             \"\"\n              '`\"``",
	}
	colors := [2]string{
		` 111111
11111  11111111111111111
     11  2      2       111
       1     2     2       11
 1     1   2    2          1 1
1 11   1                  1W1 111
 11 1111     2     1111       1111
   1     2        1  1  1     111
 11 1111   2    2  1111  111 11
1 11   1               2   11
 1   11  2    2     2   111
   111111111111111111111
   11             1111
               11111`,
		`                           111111
          11111111111111111  11111
       111       2      2  11
     11       2     2     1
    1 1          2    2   1     1
 111 1W1                  1   11 1
1111       1111     2     1111 11
 111     1  1  1        2     1
   11 111  1111  2    2   1111 11
     11   2               1   11 1
       111   2     2    2  11   1
          111111111111111111111
            1111             11
              11111`,
	}
	dir := rand.Intn(2)
	speed := 3.0
	x := -34
	if dir == 1 {
		speed = -3.0
		x = anim.Width() - 1
	}
	maxHeight := 9
	minHeight := anim.Height() - 15
	y := maxHeight
	if minHeight > maxHeight {
		y = maxHeight + rand.Intn(minHeight-maxHeight+1)
	}
	shapeFrames1 := []string{shapes[dir]}
	colorFrames1 := []string{randColor(colors[dir])}
	anim.NewEntity(NewEntityOptions{
		Shape:         shapeFrames1,
		Color:         colorFrames1,
		AutoTrans:     true,
		Position:      [3]int{x, y, Depth["shark"]},
		CallbackArgs:  []float64{speed, 0, 0},
		DeathCallback: RandomObject,
		DieOffscreen:  true,
		DefaultColor:  "YELLOW",
	})
}

func addBigFish2(_ *Entity, anim *Animation) {
	shapes := [2]string{
		"                _ _ _\n             .='\\ \\ \\`\"=,\n           .'\\ \\ \\ \\ \\ \\ \\\n\\'=._     / \\ \\ \\_\\_\\_\\_\\_\\\n\\'=._'.  /\\ \\,-\"`- _ - _ - '-.\n  \\`=._\\|'.\\/- _ - _ - _ - _- \\\n  ;\"= ._\\=./_ -_ -_ {`\"=_    @ \\\n   ;=\"_-_=- _ -  _ - {\"=_\"-     \\\n   ;_=_--_.,          {_.='   .-/\n  ;.=\"` / ';\\        _.     _.-`\n  /_.='/ \\/ /;._ _ _{.-;`/\"\n/._=_.'   '/ / / / /{.= /\n/.='       `'./_/_.=`{_/",
		"            _ _ _\n        ,=\"`/ / /'=. \n       / / / / / / /'.\n      /_/_/_/_/_/ / / \\     _.='/\n   .-' - _ - _ -`\"-,/ /\\  .'_.='/\n  / -_ - _ - _ - _ -\\/.'|/_.=`/\n / @    _=\"`} _- _- _\\.=/_. =\";\n/     -\"_=\"}  - _  - _ -=_-_\"=;\n\\-.   '=._}          ,._--_=_; \n `-._     ._        /;' \\ `\"=.;\n     `\"\\`;-.}_ _ _.;\\ \\/ \\'=._\\\n        \\ =.}\\ \\ \\ \\ \\'   '._=_.\\\n         \\_}`=._\\_\\.'`       '=.\\",
	}
	colors := [2]string{
		`                1 1 1
             1111 1 11111
           111 1 1 1 1 1 1
11111     1 1 1 11111111111
1111111  11 111112 2 2 2 2 111
  111111111112 2 2 2 2 2 2 22 1
  111 1111 12 22 22 11111    W 1
   11111112 2 2  2 2 111111     1
   111111111          11111   111
  11111 11111        11     1111
  111111 11 1111 1 111111111
1111111   11 1 1 1 1111 1
1111       1111111111111`,
		`            1 1 1
        11111 1 1111
       1 1 1 1 1 1 111
      11111111111 1 1 1     11111
   111 2 2 2 2 211111 11  1111111
  1 22 2 2 2 2 2 2 211111111111
 1 W    11111 22 22 2111111 111
1     111111 2 2  2 2 21111111
111   11111          111111111
 1111     11        111 1 11111
     111111111 1 1111 11 111111
        1 1111 1 1 1 11   1111111
         1111111111111       1111`,
	}
	dir := rand.Intn(2)
	speed := 2.5
	x := -33
	if dir == 1 {
		speed = -2.5
		x = anim.Width() - 1
	}
	maxHeight := 9
	minHeight := anim.Height() - 14
	y := maxHeight
	if minHeight > maxHeight {
		y = maxHeight + rand.Intn(minHeight-maxHeight+1)
	}
	shapeFrames2 := []string{shapes[dir]}
	colorFrames2 := []string{randColor(colors[dir])}
	anim.NewEntity(NewEntityOptions{
		Shape:         shapeFrames2,
		Color:         colorFrames2,
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
const (
	fishhookLineHeight    = 50
	fishhookTopClampY     = -10
	hookPointYOffset      = 2
	fishlineYOffsetFromHook = -fishhookLineHeight
)

func AddFishhook(_ *Entity, anim *Animation) {
	x := 10 + rand.Intn(maxInt(1, anim.Width()-30))
	yStart := -20
	yLine := yStart + fishlineYOffsetFromHook
	debugf("[fishhook] AddFishhook called: x=%d yStart=%d yLine=%d\n", x, yStart, yLine)
	anim.NewEntity(NewEntityOptions{
		EntityType:   "fishline",
		Shape:        []string{strings.Repeat("|\n", fishhookLineHeight)},
		Position:     [3]int{x + 7, yLine, Depth["water_line1"]},
		AutoTrans:    true,
		Callback:     FishhookCallback,
		CallbackArgs: map[string]string{"mode": "lowering"},
	})
	anim.NewEntity(NewEntityOptions{
		EntityType:    "fishhook",
		Shape:         []string{"       o\n      ||\n      ||\n/ \\   ||\n  \\__//\n  `--'"},
		Position:      [3]int{x, yStart, Depth["water_line1"]},
		AutoTrans:     true,
		DieOffscreen:  false,
		DefaultColor:  "GREEN",
		Callback:      FishhookCallback,
		CallbackArgs:  map[string]string{"mode": "lowering"},
		DeathCallback: func(ent *Entity, a *Animation) { GroupDeath(ent, a, []string{"hook_point", "fishline"}) },
	})
	anim.NewEntity(NewEntityOptions{
		EntityType:   "hook_point",
		Shape:        []string{".\n \n\\\n "},
		Position:     [3]int{x + 1, yStart + hookPointYOffset, Depth["shark"] + 1},
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
	nextHookY := func(y float64, hookMode string) float64 {
		if hookMode == "hooked" {
			y -= 2
			if y < float64(fishhookTopClampY) {
				y = float64(fishhookTopClampY)
			}
			return y
		}
		maxDepth := float64(int(float64(anim.Height()) * 0.75))
		y += 2
		if y > maxDepth {
			y = maxDepth
		}
		return y
	}

	mode := ""
	switch args := entity.CallbackArgs.(type) {
	case map[string]string:
		mode = args["mode"]
	}

	// Keep rig parts (fishline, hook_point) anchored to the hook while it moves.
	// Caught fish use FishhookCallback too, but they rise independently — no lookup needed.
	if entity.EntityType == "fishline" || entity.EntityType == "hook_point" {
		hooks := anim.GetEntitiesByType("fishhook")
		if len(hooks) == 0 {
			debugf("[fishhook] WARNING: %s has no parent fishhook (orphaned rig part)\n", entity.EntityType)
		} else {
			hook := hooks[0]
			hookMode := ""
			switch args := hook.CallbackArgs.(type) {
			case map[string]string:
				hookMode = args["mode"]
			}
			targetHookY := nextHookY(hook.Y, hookMode)
			switch entity.EntityType {
			case "fishline":
				entity.Y = targetHookY + float64(fishlineYOffsetFromHook)
				return true
			case "hook_point":
				entity.Y = targetHookY + float64(hookPointYOffset)
				return true
			}
		}
	}

	prevY := entity.Y
	if mode == "hooked" {
		entity.Y = nextHookY(entity.Y, mode)
	} else {
		entity.Y = nextHookY(entity.Y, mode)
	}
	maxDepth := float64(int(float64(anim.Height()) * 0.75))
	if prevY < 0 && entity.Y >= 0 {
		debugf("[fishhook] fishhook entered screen: Y=%.0f mode=%s\n", entity.Y, mode)
	}
	if prevY < maxDepth && entity.Y >= maxDepth {
		debugf("[fishhook] fishhook reached max depth: Y=%.0f maxDepth=%.0f mode=%s DieOffscreen=%v\n", entity.Y, maxDepth, mode, entity.DieOffscreen)
	}
	if entity.Y <= float64(fishhookTopClampY) && mode == "hooked" {
		debugf("[fishhook] fishhook clamped at top: Y=%.0f DieOffscreen=%v\n", entity.Y, entity.DieOffscreen)
	}
	return true
}

// Retract switches an entity into "hooked" upward movement.
// Used for fish, line, and hook after a catch event.
func Retract(entity *Entity, _ *Animation) {
	debugf("[fishhook] Retract called on %s at Y=%.0f\n", entity.EntityType, entity.Y)
	entity.Physical = false
	if entity.EntityType == "fish" {
		entity.Z = float64(Depth["water_gap2"])
		entity.Callback = FishhookCallback
		entity.CallbackArgs = map[string]string{"mode": "hooked"}
		return
	}
	entity.CallbackArgs = map[string]string{"mode": "hooked"}
	if entity.EntityType == "fishhook" {
		entity.DieOffscreen = true
		debugf("[fishhook] fishhook DieOffscreen enabled, will retract and die\n")
	}
}

// GroupDeath removes all entities of listed types from scene.
// It is used for grouped cleanup (for example hook + line + point).
// After cleanup, it chains into the next random event.
func GroupDeath(entity *Entity, anim *Animation, boundTypes []string) {
	debugf("[fishhook] GroupDeath: fishhook died at Y=%.0f, cleaning up %v\n", entity.Y, boundTypes)
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
		Color:         []string{colors[dir]},
		AutoTrans:     true,
		Position:      [3]int{x, 5, Depth["water_gap3"]},
		CallbackArgs:  []float64{speed, 0, 0, 0.25},
		DeathCallback: RandomObject,
		DieOffscreen:  true,
		DefaultColor:  "WHITE",
	})
}

// dolphinDelayOffscreenDeath runs default movement but keeps DieOffscreen false
// until the sprite overlaps the drawable area once, so formation members that
// start fully off-screen are not removed before they enter view.
func dolphinDelayOffscreenDeath(e *Entity, anim *Animation) bool {
	moved := e.MoveEntity(anim)
	if e.DieOffscreen {
		return moved
	}
	sw := float64(anim.Width())
	sh := float64(anim.Height())
	if e.X+float64(e.width) >= 0 && e.X < sw && e.Y+float64(e.height) >= 0 && e.Y < sh {
		e.DieOffscreen = true
	}
	return moved
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
		defaultColor := "CYAN"
		if i == 0 {
			defaultColor = "BLUE"
		} else if i == 1 {
			defaultColor = "MAGENTA"
		}
		if i == 0 {
			deathCb = func(_ *Entity, a *Animation) { RandomObject(nil, a) }
		}
		anim.NewEntity(NewEntityOptions{
			EntityType: "dolphin",
			Shape:      shapes[dir],
			Color:      []string{colors[dir]},
			AutoTrans:  true,
			Position:   [3]int{x - (distance * (2 - i)), 5, Depth["water_gap3"]},
			Callback:   dolphinDelayOffscreenDeath,
			CallbackArgs: []float64{
				speed, 0, 0, 0.5,
			},
			DeathCallback: deathCb,
			DieOffscreen:  false,
			DefaultColor:  defaultColor,
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
		Shape:         []string{shapes[dir]},
		Color:         []string{colors[dir]},
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
	names := []string{
		"AddShip", "AddWhale", "AddMonster", "AddBigFish", "AddShark",
		"AddFishhook", "AddSwan", "AddDucks", "AddDolphins",
	}
	idx := rand.Intn(len(randomObjects))
	debugf("[RandomObject] picked %s (index %d)\n", names[idx], idx)
	randomObjects[idx](dead, anim)
}

// maxInt returns the larger of two integers.
// It is a small helper for safe random ranges.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
