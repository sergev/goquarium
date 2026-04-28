package aquarium

import (
	"testing"
	"time"
)

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

func TestEntityShouldDie(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Second)
	e := NewEntity(NewEntityOptions{Shape: "abc", DieTime: &past})
	if !e.ShouldDie(100, 100, now) {
		t.Fatalf("expected entity to die by time")
	}
}

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
