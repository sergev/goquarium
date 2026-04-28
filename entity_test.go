package main

import (
	"math"
	"strings"
	"testing"
	"time"
)

// TestEntityMoveFrame checks default movement and frame progression.
// Two updates are used because 0.6 + 0.6 crosses the frame threshold.
// This protects the "frameStep >= 1" animation rule.
func TestEntityMoveFrame(t *testing.T) {
	e := NewEntity(NewEntityOptions{
		Shape:        []string{"a", "b"},
		CallbackArgs: []float64{1, 2, 0, 0.6},
	})
	e.MoveEntity(nil)
	if e.X != 1 || e.Y != 2 {
		t.Fatalf("unexpected position: %v,%v", e.X, e.Y)
	}
	e.MoveEntity(nil)
	if e.CurrentFrame == 0 {
		t.Fatalf("expected frame to advance")
	}
}

// TestEntityShouldDie verifies time-based death rule.
// Entity should be removed when current time is past die time.
func TestEntityShouldDie(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Second)
	e := NewEntity(NewEntityOptions{Shape: "abc", DieTime: &past})
	if !e.ShouldDie(100, 100, now) {
		t.Fatalf("expected entity to die by time")
	}
}

// TestCollisionDetection verifies rectangle overlap collision behavior.
// This test protects the simple AABB rule used in checkCollisions.
func TestCollisionDetection(t *testing.T) {
	a := NewAnimation()
	e1 := NewEntity(NewEntityOptions{Shape: "xx", Position: [3]int{1, 1, 1}, Physical: true})
	e2 := NewEntity(NewEntityOptions{Shape: "xx", Position: [3]int{2, 1, 1}})
	a.entities = []*Entity{e1, e2}
	a.checkCollisions()
	if len(e1.Collision) == 0 {
		t.Fatalf("expected collision")
	}
}

// TestRandColor ensures color placeholders are replaced.
// Output should contain only known color marker runes.
func TestRandColor(t *testing.T) {
	mask := randColor("123456789")
	if len(mask) != 9 {
		t.Fatalf("invalid transformed mask")
	}
	for _, ch := range mask {
		switch ch {
		case 'c', 'C', 'r', 'R', 'y', 'Y', 'b', 'B', 'g', 'G', 'm', 'M':
		default:
			t.Fatalf("unexpected color rune %q", ch)
		}
	}
}

// TestSurfaceSpritesUseRealMultilineStrings protects sprite text integrity.
// Literal "\\n" would flatten art into one broken line on screen.
// This regression test keeps surface animations readable.
func TestSurfaceSpritesUseRealMultilineStrings(t *testing.T) {
	anim := NewAnimation()
	anim.width = 120
	anim.height = 40

	AddShip(nil, anim)
	AddWhale(nil, anim)
	AddDucks(nil, anim)
	AddDolphins(nil, anim)
	AddSwan(nil, anim)

	for _, e := range anim.entities {
		if e.Z != float64(Depth["water_gap1"]) &&
			e.Z != float64(Depth["water_gap2"]) &&
			e.Z != float64(Depth["water_gap3"]) {
			continue
		}

		shape := e.CurrentShape()
		if !strings.Contains(shape, "\n") {
			t.Fatalf("expected multiline surface sprite for %q", e.EntityType)
		}
		if strings.Contains(shape, `\n`) {
			t.Fatalf("found literal \\n sequence in surface sprite for %q", e.EntityType)
		}
	}
}

// TestSeaMonsterSpritesParity protects restored sea monster behavior.
// It ensures monsters are animated, multiline, color-masked, and move at expected speed.
func TestSeaMonsterSpritesParity(t *testing.T) {
	anim := NewAnimation()
	anim.width = 160
	anim.height = 50

	for i := 0; i < 40; i++ {
		before := len(anim.entities)
		AddMonster(nil, anim)
		if len(anim.entities) != before+1 {
			t.Fatalf("expected AddMonster to add one entity, got before=%d after=%d", before, len(anim.entities))
		}

		e := anim.entities[len(anim.entities)-1]
		if len(e.Shapes) < 2 {
			t.Fatalf("expected animated monster with at least 2 frames, got %d", len(e.Shapes))
		}
		for _, frame := range e.Shapes {
			if !strings.Contains(frame, "\n") {
				t.Fatalf("expected multiline monster frame")
			}
			if strings.Contains(frame, `\n`) {
				t.Fatalf("found literal \\n sequence in monster frame")
			}
		}

		if len(e.Colors) == 0 || e.Colors[0] == "" {
			t.Fatalf("expected monster color masks")
		}
		if len(e.Colors) != len(e.Shapes) {
			t.Fatalf("expected color mask frame count to match shape frame count, got %d vs %d", len(e.Colors), len(e.Shapes))
		}

		args, ok := e.CallbackArgs.([]float64)
		if !ok {
			t.Fatalf("expected []float64 callback args, got %T", e.CallbackArgs)
		}
		if len(args) != 4 {
			t.Fatalf("expected 4 callback args, got %d", len(args))
		}
		if math.Abs(args[0]) != 2.0 || args[1] != 0 || args[2] != 0 || args[3] != 0.25 {
			t.Fatalf("unexpected monster callback args: %#v", args)
		}
	}
}

func TestFishDesignCatalogParityCounts(t *testing.T) {
	if got := len(oldFishDesigns); got != 8 {
		t.Fatalf("expected 8 old fish designs, got %d", got)
	}
	if got := len(newFishDesigns); got != 4 {
		t.Fatalf("expected 4 new fish designs, got %d", got)
	}
}

func TestBigFishVisualParity(t *testing.T) {
	anim := NewAnimation()
	anim.width = 200
	anim.height = 60

	for i := 0; i < 30; i++ {
		before := len(anim.entities)
		AddBigFish(nil, anim)
		if len(anim.entities) != before+1 {
			t.Fatalf("expected one new big fish entity")
		}
		e := anim.entities[len(anim.entities)-1]
		shape := e.CurrentShape()
		if !strings.Contains(shape, "\n") || strings.Contains(shape, `\n`) {
			t.Fatalf("big fish shape should be real multiline sprite")
		}
		if len(strings.Split(shape, "\n")) < 10 {
			t.Fatalf("expected large big fish sprite, got too few lines")
		}
		color := e.CurrentColor()
		if color == "" {
			t.Fatalf("expected big fish color mask")
		}
		args, ok := e.CallbackArgs.([]float64)
		if !ok || len(args) < 1 {
			t.Fatalf("unexpected callback args for big fish: %T", e.CallbackArgs)
		}
		if math.Abs(args[0]) != 2.5 && math.Abs(args[0]) != 3.0 {
			t.Fatalf("unexpected big fish speed: %v", args[0])
		}
	}
}

func TestSharkAndTeethParity(t *testing.T) {
	anim := NewAnimation()
	anim.width = 200
	anim.height = 60
	AddShark(nil, anim)

	var shark *Entity
	var teeth *Entity
	for _, e := range anim.entities {
		if e.EntityType == "shark" {
			shark = e
		}
		if e.EntityType == "teeth" {
			teeth = e
		}
	}
	if shark == nil || teeth == nil {
		t.Fatalf("expected both shark and teeth entities")
	}
	shape := shark.CurrentShape()
	if !strings.Contains(shape, "\n") || strings.Contains(shape, `\n`) {
		t.Fatalf("shark shape should be real multiline sprite")
	}
	if len(strings.Split(shape, "\n")) < 8 {
		t.Fatalf("expected large shark sprite, got too few lines")
	}
	if shark.CurrentColor() == "" {
		t.Fatalf("expected shark color mask")
	}
	sx, sy, _ := shark.Position()
	tx, ty, _ := teeth.Position()
	if ty != sy+7 {
		t.Fatalf("expected teeth y offset +7, got shark=%d teeth=%d", sy, ty)
	}
	if tx != -9 && tx != sx+9 {
		t.Fatalf("unexpected teeth x alignment: shark=%d teeth=%d", sx, tx)
	}
}

func TestFishhookVisualParity(t *testing.T) {
	anim := NewAnimation()
	anim.width = 200
	anim.height = 60
	AddFishhook(nil, anim)

	var hook *Entity
	var line *Entity
	for _, e := range anim.entities {
		if e.EntityType == "fishhook" {
			hook = e
		}
		if e.EntityType == "fishline" {
			line = e
		}
	}
	if hook == nil || line == nil {
		t.Fatalf("expected fishhook and fishline entities")
	}
	if len(strings.Split(hook.CurrentShape(), "\n")) < 6 {
		t.Fatalf("expected larger fishhook sprite")
	}
	if strings.Count(line.CurrentShape(), "|\n") < 50 {
		t.Fatalf("expected long fishline with at least 50 pipe segments")
	}
}
