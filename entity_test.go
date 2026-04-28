package main

import (
	"strings"
	"testing"
	"time"
)

// TestEntityMoveFrame checks default movement math.
// It also confirms frame index advances after enough frame time.
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

// TestCollisionDetection checks overlap collision logic.
// Two close entities should appear in collision list.
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

// TestSurfaceSpritesUseRealMultilineStrings guards sprite formatting.
// It catches accidental literal "\\n" strings in surface sprites.
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
