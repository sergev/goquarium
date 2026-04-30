package main

import (
	"math/rand"
	"strings"
)

// AddBubble creates one rising bubble from a fish.
// CallbackArgs follow [dx, dy, dz, frameStep], so dy=-1 means "go up".
// Bubble starts on fish mouth side based on fish horizontal direction.
func AddBubble(fish *Entity, anim *Animation) {
	cbArgs, _ := fish.CallbackArgs.([]float64)
	fw, fh := fish.Size()
	fx, fy, fz := fish.Position()
	bx := fx
	if len(cbArgs) > 0 && cbArgs[0] > 0 {
		bx += fw
	}
	anim.NewEntity(NewEntityOptions{
		EntityType:   "bubble",
		Shape:        []string{".", "o", "O", "O", "O"},
		Position:     [3]int{bx, fy + fh/2, fz - 1},
		CallbackArgs: []float64{0, -1, 0, 0.1},
		DieOffscreen: true,
		Physical:     true,
		CollHandler:  BubbleCollision,
		DefaultColor: "CYAN",
	})
}

// BubbleCollision removes a bubble when it reaches water line.
// This keeps bubbles from floating forever.
func BubbleCollision(bubble *Entity, _ *Animation) {
	for _, obj := range bubble.Collision {
		if obj.EntityType == "waterline" {
			bubble.Kill()
			return
		}
	}
}

// FishCallback adds random bubble behavior to fish.
// After that, it uses normal movement logic.
func FishCallback(fish *Entity, anim *Animation) bool {
	if rand.Intn(100)+1 > 97 {
		AddBubble(fish, anim)
	}
	return fish.MoveEntity(anim)
}

// FishCollision handles dangerous contacts for a fish.
// Shark teeth can kill small fish; hook contact retracts multiple entities.
// One collision may change state of fish, hook, line, and hook point together.
func FishCollision(fish *Entity, anim *Animation) {
	for _, obj := range fish.Collision {
		if obj.EntityType == "teeth" {
			_, h := fish.Size()
			if h <= 5 {
				x, y, z := obj.Position()
				AddSplat(anim, x, y, z)
				fish.Kill()
			}
			return
		}
		if obj.EntityType == "hook_point" && obj.Physical {
			Retract(obj, anim)
			Retract(fish, anim)
			for _, h := range anim.GetEntitiesByType("fishhook") {
				Retract(h, anim)
			}
			for _, l := range anim.GetEntitiesByType("fishline") {
				Retract(l, anim)
			}
			return
		}
	}
}

// AddSplat creates a short visual burst when fish gets eaten.
// DieFrame means "remove after enough frame steps," not wall-clock seconds.
func AddSplat(anim *Animation, x, y, z int) {
	frames := []string{
		"\n\n   .\n  ***\n   '\n\n",
		"\n\n .,*;`\n '*,**\n *'~'\n\n",
		"\n  , ,\n \" ,\"'\n *\" *'\"\n  \" ; .\n\n",
		"* ' , ' `\n' ` * . '\n ' `' \",'\n* ' \" * .\n\" * ', '",
	}
	anim.NewEntity(NewEntityOptions{
		Shape:        frames,
		Position:     [3]int{x - 4, y - 2, z - 2},
		DefaultColor: "RED",
		CallbackArgs: []float64{0, 0, 0, 0.25},
		AutoTrans:    true,
		DieFrame:     15,
	})
}

// fishDesign stores left/right sprite variants and color masks.
// We pick one direction based on travel direction.
type fishDesign struct {
	shape [2]string
	color [2]string
}

var oldFishDesigns = []fishDesign{
	{shape: [2]string{
		"       \\\n     ...\\..,\n\\  /'       \\\n >=     (  ' >\n/  \\      / /\n    `\"'\"'/''",
		"      /\n  ,../...\n /       '\\  /\n< '  )     =<\n \\ \\      /  \\\n  `'\"'\"'",
	}, color: [2]string{
		"       2\n     1112111\n6  11       1\n 66     7  4 5\n6  1      3 1\n    11111311",
		"      2\n  1112111\n 1       11  6\n5 4  7     66\n 1 3      1  6\n  11311111",
	}},
	{shape: [2]string{
		"    \\\n\\ /--\\\n>=  (o>\n/ \\__/\n    /",
		"  /\n /--\\ /\n<o)  =<\n \\__/ \\\n  \\",
	}, color: [2]string{
		"    2\n6 1111\n66  745\n6 1111\n    3",
		"  2\n 1111 6\n547  66\n 1111 6\n  3",
	}},
	{shape: [2]string{
		"       \\:.\n\\;,   ,;\\\\\\\\,,\n  \\\\\\\\;;:::::::o\n  ///;;::::::::<\n /;` ``/////``",
		"      .:/\n   ,,///;,   ,;/\n o:::::::;;///\n>::::::::;;\\\\\\\\\n  ''\\\\\\\\\\\\\\\\'' ';\\",
	}, color: [2]string{
		"       222\n666   1122211\n  6661111111114\n  66611111111115\n 666 113333311",
		"      222\n   1122211   666\n 4111111111666\n51111111111666\n  113333311 666",
	}},
	{shape: [2]string{
		"  __\n><_'>\n   '",
		" __\n<'_><\n `",
	}, color: [2]string{
		"  11\n61145\n   3",
		" 11\n54116\n 3",
	}},
	{shape: [2]string{
		"   ..\\\\\n>='   ('>\n  '''/''",
		"  ,..\n<')   `=<\n ``\\```",
	}, color: [2]string{
		"   1121\n661   745\n  111311",
		"  1211\n547   166\n 113111",
	}},
	{shape: [2]string{
		"   \\\n  / \\\n>=_('>\n  \\_/\n   /",
		"  /\n / \\\n<')_=<\n \\_/\n  \\",
	}, color: [2]string{
		"   2\n  1 1\n661745\n  111\n   3",
		"  2\n 1 1\n547166\n 111\n  3",
	}},
	{shape: [2]string{
		"  ,\\\n>=('>\n  '/",
		" /,\n<')=<\n \\`",
	}, color: [2]string{
		"  12\n66745\n  13",
		" 21\n54766\n 31",
	}},
	{shape: [2]string{
		"  __\n\\/ o\\\n/\\__/",
		" __\n/o \\/\n\\__/\\",
	}, color: [2]string{
		"  11\n61 41\n61111",
		" 11\n14 16\n11116",
	}},
}

var newFishDesigns = []fishDesign{
	{shape: [2]string{
		"   \\\n  / \\\n>=_('>\n  \\_/\n   /",
		"  /\n / \\\n<')_=<\n \\_/\n  \\",
	}, color: [2]string{
		"   1\n  1 1\n663745\n  111\n   3",
		"  2\n 111\n547366\n 111\n  3",
	}},
	{shape: [2]string{
		"     ,\n     }\\\\\n\\  .'  `\\\n}}<   ( 6>\n/  `,  .'\n     }/\n     '",
		"    ,\n   /{\n /'  `.  /\n<6 )   >{{\n `.  ,'  \\\n   {\\\n    `",
	}, color: [2]string{
		"     2\n     22\n6  11  11\n661   7 45\n6  11  11\n     33\n     3",
		"    2\n   22\n 11  11  6\n54 7   166\n 11  11  6\n   33\n    3",
	}},
	{shape: [2]string{
		"            \\'`.\n             )  \\\n(`.      _.-`' ' '`-.\n \\ `.  .`        (o) \\_\n  >  ><     (((       (\n / .`  ._      /_|  /'\n(.`       `-. _  _.-`\n            /__/'",
		"       .'`/\n      /  (\n  .-'` ` `'-._      .')\n_/ (o)        '.  .' /\n)       )))     ><  <\n`\\  |_\\      _.'  '. \\\n  '-._  _ .-'       '.)\n      `\\__\\",
	}, color: [2]string{
		"            1111\n             1  1\n111      11111 1 1111\n 1 11  11        141 11\n  1  11     777       5\n 1 11  111      333  11\n111       111 1  1111\n            11111",
		"       1111\n      1  1\n  1111 1 11111      111\n11 141        11  11 1\n5       777     11  1\n11  333      111  11 1\n  1111  1 111       111\n      11111",
	}},
	{shape: [2]string{
		"       ,--,_\n__    _\\.---'-.\n\\ '.-\"     // o\\\n/_.'-._    \\\\  /\n       `\"--(/\"`",
		"    _,--,\n .-'---./_    __\n/o \\\\     \"-.' /\n\\  //    _.-'._\\\n `\"\\)--\"`",
	}, color: [2]string{
		"       22222\n66    121111211\n6 6111     77 41\n6661111    77  1\n       11113311",
		"    22222\n 112111121    66\n14 77     1116 6\n1  77    1111666\n 11331111",
	}},
}

// randColor replaces numeric color placeholders in sprite masks.
// Digits 1..9 are templates that become random color marker letters.
// This keeps fish colors varied without changing shape art.
func randColor(mask string) string {
	colors := []string{"c", "C", "r", "R", "y", "Y", "b", "B", "g", "G", "m", "M"}
	out := mask
	for i := 1; i <= 9; i++ {
		out = strings.ReplaceAll(out, string(rune('0'+i)), colors[rand.Intn(len(colors))])
	}
	return out
}

// AddFish spawns one fish using mode rules and random design.
// Direction chooses speed sign and which side of screen fish starts from.
// Death callback respawns another fish, keeping population stable.
func AddFish(_ *Entity, anim *Animation, classicMode bool) {
	var design fishDesign
	if classicMode || rand.Intn(12)+1 <= 8 {
		design = oldFishDesigns[rand.Intn(len(oldFishDesigns))]
	} else {
		design = newFishDesigns[rand.Intn(len(newFishDesigns))]
	}
	direction := rand.Intn(2)
	speed := 0.25 + rand.Float64()*1.75
	if direction == 1 {
		speed *= -1
	}
	depth := Depth["fish_start"] + rand.Intn(Depth["fish_end"]-Depth["fish_start"]+1)
	shapeFrames := []string{design.shape[direction]}
	colorFrames := []string{randColor(design.color[direction])}
	fish := NewEntity(NewEntityOptions{
		EntityType:    "fish",
		Shape:         shapeFrames,
		Color:         colorFrames,
		AutoTrans:     true,
		Position:      [3]int{0, 0, depth},
		Callback:      FishCallback,
		CallbackArgs:  []float64{speed, 0, 0},
		DieOffscreen:  true,
		Physical:      true,
		CollHandler:   FishCollision,
		DeathCallback: func(ent *Entity, a *Animation) { AddFish(ent, a, classicMode) },
	})

	waterBottom := 9
	screenBottom := anim.Height() - 1
	_, fh := fish.Size()
	available := screenBottom - waterBottom - fh
	if available > 0 {
		fish.Y = float64(waterBottom + rand.Intn(available+1))
	} else {
		fish.Y = float64(waterBottom)
	}
	fw, _ := fish.Size()
	if direction == 0 {
		fish.X = float64(-fw)
	} else {
		fish.X = float64(anim.Width())
	}
	anim.AddEntity(fish)
}

// AddAllFish creates initial fish count from screen area.
// The /350 constant is a simple density tuning value.
func AddAllFish(anim *Animation, classicMode bool) {
	screenSize := (anim.Height() - 9) * anim.Width()
	count := screenSize / 350
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		AddFish(nil, anim, classicMode)
	}
}
