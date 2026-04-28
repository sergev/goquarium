package aquarium

import (
	"math/rand"
	"strings"
)

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

func BubbleCollision(bubble *Entity, _ *Animation) {
	for _, obj := range bubble.Collision {
		if obj.EntityType == "waterline" {
			bubble.Kill()
			return
		}
	}
}

func FishCallback(fish *Entity, anim *Animation) bool {
	if rand.Intn(100)+1 > 97 {
		AddBubble(fish, anim)
	}
	return fish.MoveEntity(anim)
}

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
		if obj.EntityType == "hook_point" {
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

type fishDesign struct {
	shape [2]string
	color [2]string
}

var oldFishDesigns = []fishDesign{
	{shape: [2]string{
		"  __\n><_'>\n   '",
		" __\n<'_><\n `",
	}, color: [2]string{
		"  11\n61145\n   3",
		" 11\n54116\n 3",
	}},
	{shape: [2]string{
		"  ,\\\n>=('>\n  '/",
		" /,\n<')=<\n \\`",
	}, color: [2]string{
		"  12\n66745\n  13",
		" 21\n54766\n 31",
	}},
	{shape: [2]string{
		"   \\\n  / \\\n>=_('>\n  \\_/\n   /",
		"  /\n / \\\n<')_=<\n \\_/\n  \\",
	}, color: [2]string{
		"   2\n  1 1\n661745\n  111\n   3",
		"  2\n 1 1\n547166\n 111\n  3",
	}},
}

var newFishDesigns = []fishDesign{
	{shape: [2]string{
		"     ,\n     }\\\\\n\\  .'  `\\\n}}<   ( 6>\n/  `,  .'\n     }/\n     '",
		"    ,\n   /{\n /'  `.  /\n<6 )   >{{\n `.  ,'  \\\n   {\\\n    `",
	}, color: [2]string{
		"     2\n     22\n6  11  11\n661   7 45\n6  11  11\n     33\n     3",
		"    2\n   22\n 11  11  6\n54 7   166\n 11  11  6\n   33\n    3",
	}},
}

func randColor(mask string) string {
	colors := []string{"c", "C", "r", "R", "y", "Y", "b", "B", "g", "G", "m", "M"}
	out := mask
	for i := 1; i <= 9; i++ {
		out = strings.ReplaceAll(out, string(rune('0'+i)), colors[rand.Intn(len(colors))])
	}
	return out
}

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
	fish := NewEntity(NewEntityOptions{
		EntityType:    "fish",
		Shape:         design.shape[direction],
		Color:         randColor(design.color[direction]),
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
