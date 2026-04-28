package aquarium

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Animation struct {
	screen       tcell.Screen
	entities     []*Entity
	colorEnabled bool
	running      bool
	width        int
	height       int
	maskColorMap map[rune]tcell.Color
}

func NewAnimation() *Animation {
	return &Animation{
		colorEnabled: true,
		maskColorMap: map[rune]tcell.Color{
			'r': tcell.ColorRed, 'R': tcell.ColorRed,
			'g': tcell.ColorGreen, 'G': tcell.ColorGreen,
			'y': tcell.ColorYellow, 'Y': tcell.ColorYellow,
			'b': tcell.ColorBlue, 'B': tcell.ColorBlue,
			'm': tcell.ColorDarkMagenta, 'M': tcell.ColorDarkMagenta,
			'c': tcell.ColorTeal, 'C': tcell.ColorTeal,
			'w': tcell.ColorWhite, 'W': tcell.ColorWhite,
			'k': tcell.ColorBlack, 'K': tcell.ColorBlack,
			'1': tcell.ColorTeal, '2': tcell.ColorYellow, '3': tcell.ColorGreen,
			'4': tcell.ColorWhite, '5': tcell.ColorRed, '6': tcell.ColorBlue,
			'7': tcell.ColorDarkMagenta, '8': tcell.ColorBlack, '9': tcell.ColorWhite,
		},
	}
}

func (a *Animation) Width() int  { return a.width }
func (a *Animation) Height() int { return a.height }

func (a *Animation) NewEntity(opts NewEntityOptions) *Entity {
	e := NewEntity(opts)
	a.AddEntity(e)
	return e
}

func (a *Animation) AddEntity(e *Entity) {
	a.entities = append(a.entities, e)
	sort.Slice(a.entities, func(i, j int) bool { return a.entities[i].Z < a.entities[j].Z })
}

func (a *Animation) DelEntity(e *Entity) {
	for i := range a.entities {
		if a.entities[i] == e {
			a.entities = append(a.entities[:i], a.entities[i+1:]...)
			return
		}
	}
}

func (a *Animation) RemoveAllEntities() { a.entities = nil }

func (a *Animation) GetEntitiesByType(tp string) []*Entity {
	out := make([]*Entity, 0)
	for _, e := range a.entities {
		if e.EntityType == tp {
			out = append(out, e)
		}
	}
	return out
}

func (a *Animation) updateSize() error {
	w, h := a.screen.Size()
	if h < 15 || w < 40 {
		return fmt.Errorf("terminal too small: need at least 40x15, got %dx%d", w, h)
	}
	a.width = w
	a.height = h - 1
	return nil
}

func (a *Animation) checkCollisions() {
	for _, e := range a.entities {
		e.Collision = nil
	}
	for _, e := range a.entities {
		if !e.Physical {
			continue
		}
		ex, ey, _ := e.Position()
		ew, eh := e.Size()
		for _, o := range a.entities {
			if e == o {
				continue
			}
			ox, oy, _ := o.Position()
			ow, oh := o.Size()
			if ex < ox+ow && ex+ew > ox && ey < oy+oh && ey+eh > oy {
				e.Collision = append(e.Collision, o)
			}
		}
	}
}

func colorByName(name string) tcell.Color {
	switch strings.ToUpper(name) {
	case "BLACK":
		return tcell.ColorBlack
	case "RED":
		return tcell.ColorRed
	case "GREEN":
		return tcell.ColorGreen
	case "YELLOW":
		return tcell.ColorYellow
	case "BLUE":
		return tcell.ColorBlue
	case "MAGENTA":
		return tcell.ColorDarkMagenta
	case "CYAN":
		return tcell.ColorTeal
	default:
		return tcell.ColorWhite
	}
}

func (a *Animation) drawEntity(e *Entity) {
	x, y, _ := e.Position()
	lines := strings.Split(e.CurrentShape(), "\n")
	colorLines := strings.Split(e.CurrentColor(), "\n")
	for li, line := range lines {
		drawY := y + li
		if drawY < 0 || drawY >= a.height {
			continue
		}
		lineRunes := []rune(line)
		var colorRunes []rune
		if li < len(colorLines) {
			colorRunes = []rune(colorLines[li])
		}
		for ci, ch := range lineRunes {
			drawX := x + ci
			if drawX < 0 || drawX >= a.width {
				continue
			}
			if e.AutoTrans && (ch == ' ' || ch == e.Transparent) {
				continue
			}
			if ch < 32 {
				continue
			}
			style := tcell.StyleDefault.Foreground(colorByName(e.DefaultColor))
			if a.colorEnabled && ci < len(colorRunes) {
				if c, ok := a.maskColorMap[colorRunes[ci]]; ok {
					style = style.Foreground(c)
				}
			}
			a.screen.SetContent(drawX, drawY, ch, nil, style)
		}
	}
}

func (a *Animation) animate() {
	now := time.Now()
	for _, e := range append([]*Entity{}, a.entities...) {
		e.Update(a)
	}
	a.checkCollisions()
	for _, e := range append([]*Entity{}, a.entities...) {
		if e.ShouldDie(a.width, a.height, now) {
			if e.DeathCallback != nil {
				e.DeathCallback(e, a)
			}
			a.DelEntity(e)
		}
	}
	a.screen.Clear()
	sorted := append([]*Entity{}, a.entities...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Z > sorted[j].Z })
	for _, e := range sorted {
		a.drawEntity(e)
	}
	a.screen.Show()
}

func (a *Animation) drawInfoOverlay() {
	lines := InfoLines()
	a.screen.Clear()
	startY := (a.height - len(lines)) / 2
	if startY < 0 {
		startY = 0
	}
	for i, ln := range lines {
		y := startY + i
		if y >= a.height {
			break
		}
		x := (a.width - len([]rune(ln))) / 2
		if x < 0 {
			x = 0
		}
		for ci, ch := range []rune(ln) {
			if x+ci >= a.width {
				break
			}
			a.screen.SetContent(x+ci, y, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}
	}
	a.screen.Show()
}

func (a *Animation) Run(setup func(*Animation, bool), classic bool) error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := s.Init(); err != nil {
		return err
	}
	defer s.Fini()
	a.screen = s
	if err := a.updateSize(); err != nil {
		return err
	}
	a.running = true
	setup(a, classic)

	eventCh := make(chan tcell.Event, 32)
	go func() {
		for a.running {
			eventCh <- a.screen.PollEvent()
		}
	}()
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	paused := false
	showingInfo := false

	for a.running {
		select {
		case ev := <-eventCh:
			switch tev := ev.(type) {
			case *tcell.EventResize:
				a.screen.Sync()
				if err := a.updateSize(); err != nil {
					return err
				}
				if showingInfo {
					a.drawInfoOverlay()
				}
			case *tcell.EventKey:
				if tev.Key() == tcell.KeyEscape && showingInfo {
					showingInfo = false
					paused = false
					continue
				}
				if tev.Key() == tcell.KeyRune {
					switch tev.Rune() {
					case 'q', 'Q':
						a.running = false
					case 'r', 'R':
						a.RemoveAllEntities()
						setup(a, classic)
					case 'p', 'P':
						if !showingInfo {
							paused = !paused
						}
					case 'i', 'I':
						showingInfo = !showingInfo
						if showingInfo {
							paused = true
							a.drawInfoOverlay()
						} else {
							paused = false
						}
					}
				}
			}
		case <-tick.C:
			if showingInfo {
				a.drawInfoOverlay()
				continue
			}
			if !paused {
				a.animate()
			}
		}
	}
	return nil
}

func EnsureScreenSupport() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if s == nil {
		return errors.New("no compatible terminal screen available")
	}
	return nil
}
