package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

// Animation is the main runtime controller for the aquarium.
// It owns screen state, active entities, and loop flags.
// Most game flow starts from this structure.
type Animation struct {
	entities     []*Entity
	colorEnabled bool
	running      bool
	width        int
	height       int
	maskColorMap map[rune]termbox.Attribute
}

// NewAnimation creates a new animation manager with defaults.
// It also prepares the lookup map for masked colors.
// Call this before spawning anything.
func NewAnimation() *Animation {
	return &Animation{
		colorEnabled: true,
		maskColorMap: map[rune]termbox.Attribute{
			'r': termbox.ColorRed, 'R': termbox.ColorRed,
			'g': termbox.ColorGreen, 'G': termbox.ColorGreen,
			'y': termbox.ColorYellow, 'Y': termbox.ColorYellow,
			'b': termbox.ColorBlue, 'B': termbox.ColorBlue,
			'm': termbox.ColorMagenta, 'M': termbox.ColorMagenta,
			'c': termbox.ColorCyan, 'C': termbox.ColorCyan,
			'w': termbox.ColorWhite, 'W': termbox.ColorWhite,
			'k': termbox.ColorBlack, 'K': termbox.ColorBlack,
			'1': termbox.ColorCyan, '2': termbox.ColorYellow, '3': termbox.ColorGreen,
			'4': termbox.ColorWhite, '5': termbox.ColorRed, '6': termbox.ColorBlue,
			'7': termbox.ColorMagenta, '8': termbox.ColorBlack, '9': termbox.ColorWhite,
		},
	}
}

// Width returns current drawable screen width in cells.
// Spawners use this value to place entities safely.
func (a *Animation) Width() int { return a.width }

// Height returns current drawable screen height in cells.
// This updates after terminal resize events.
func (a *Animation) Height() int { return a.height }

// NewEntity builds an entity and adds it to the world.
// This helper saves you from calling two functions manually.
func (a *Animation) NewEntity(opts NewEntityOptions) *Entity {
	e := NewEntity(opts)
	a.AddEntity(e)
	return e
}

// AddEntity adds one entity to the internal list.
// It keeps depth order stable for later drawing.
func (a *Animation) AddEntity(e *Entity) {
	a.entities = append(a.entities, e)
	sort.Slice(a.entities, func(i, j int) bool { return a.entities[i].Z < a.entities[j].Z })
}

// DelEntity removes one matching entity from the list.
// If the entity is missing, this function just returns.
func (a *Animation) DelEntity(e *Entity) {
	for i := range a.entities {
		if a.entities[i] == e {
			a.entities = append(a.entities[:i], a.entities[i+1:]...)
			return
		}
	}
}

// RemoveAllEntities clears every object from the scene.
// This is used by the reset command.
func (a *Animation) RemoveAllEntities() { a.entities = nil }

// GetEntitiesByType returns objects with the same type label.
// Behavior code uses it to find hooks, teeth, lines, etc.
func (a *Animation) GetEntitiesByType(tp string) []*Entity {
	out := make([]*Entity, 0)
	for _, e := range a.entities {
		if e.EntityType == tp {
			out = append(out, e)
		}
	}
	return out
}

// updateSize refreshes width/height after startup or terminal resize.
// We require a minimum size so large ASCII art does not break badly.
// Height is stored as one row less to avoid bottom-row terminal glitches.
func (a *Animation) updateSize() error {
	w, h := termbox.Size()
	if h < 15 || w < 40 {
		return fmt.Errorf("terminal too small: need at least 40x15, got %dx%d", w, h)
	}
	a.width = w
	a.height = h - 1
	return nil
}

// checkCollisions finds overlaps using simple rectangle checks.
// This is an O(n^2) pass, but it is easy to understand and fine here.
// Results are saved on each entity for later collision handlers.
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

// colorByName converts a color name string to a termbox attribute.
// Unknown names are treated as white so drawing still works.
func colorByName(name string) termbox.Attribute {
	switch strings.ToUpper(name) {
	case "BLACK":
		return termbox.ColorBlack
	case "RED":
		return termbox.ColorRed
	case "GREEN":
		return termbox.ColorGreen
	case "YELLOW":
		return termbox.ColorYellow
	case "BLUE":
		return termbox.ColorBlue
	case "MAGENTA":
		return termbox.ColorMagenta
	case "CYAN":
		return termbox.ColorCyan
	case "DARK_GREY":
		return termbox.ColorDarkGray
	default:
		return termbox.ColorWhite
	}
}

// drawEntity paints one object frame onto the screen grid.
// Shape and color masks are read line-by-line in parallel.
// Mask letters/digits pick colors, while transparent cells are skipped.
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
			fg := colorByName(e.DefaultColor)
			if a.colorEnabled && ci < len(colorRunes) {
				if c, ok := a.maskColorMap[colorRunes[ci]]; ok {
					fg = c
				}
			}
			termbox.SetCell(drawX, drawY, ch, fg, termbox.ColorDefault)
		}
	}
}

// drawFrame renders current entities without advancing simulation state.
// This is used for immediate redraw requests like terminal resize.
func (a *Animation) drawFrame() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	sorted := append([]*Entity{}, a.entities...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Z > sorted[j].Z })
	for _, e := range sorted {
		a.drawEntity(e)
	}
	termbox.Flush()
}

// reflowForResize rebuilds size-dependent static scenery and keeps dynamic entities.
// It also clamps preserved entity Y positions into the drawable vertical range.
func (a *Animation) reflowForResize() {
	preserved := make([]*Entity, 0, len(a.entities))
	for _, e := range a.entities {
		switch e.EntityType {
		case "waterline", "castle", "seaweed":
			continue
		default:
			_, eh := e.Size()
			maxY := a.height - eh
			if maxY < 0 {
				maxY = 0
			}
			if e.Y < 0 {
				e.Y = 0
			}
			if e.Y > float64(maxY) {
				e.Y = float64(maxY)
			}
			preserved = append(preserved, e)
		}
	}
	a.entities = preserved
	AddEnvironment(a)
	AddCastle(a)
	AddAllSeaweed(a)
}

// animate runs one full simulation step and then draws the frame.
// We copy entity slices before loops so callbacks can add/remove safely.
// Pass order is: update -> collisions -> death cleanup -> render.
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
	a.drawFrame()
}

// infoStyleFor returns a foreground attribute for one info-overlay rune.
// It keeps colors consistent across header, controls, and hint lines.
func infoStyleFor(lineIdx int, line string, ch rune) termbox.Attribute {
	base := termbox.ColorWhite
	if ch == ' ' {
		return base
	}

	// Box frame and header lines use an accent color.
	if lineIdx <= 4 {
		if strings.ContainsRune("╔═╗║╚╝", ch) {
			return termbox.ColorCyan
		}
		return termbox.ColorWhite
	}

	// Controls line highlights control keys and ESC.
	if strings.Contains(line, "Q/q quit") {
		return termbox.ColorGreen
	}

	// Final hint line uses a secondary accent and emphasizes key names.
	if strings.Contains(line, "Press I or ESC") {
		return termbox.ColorMagenta
	}

	return base
}

// drawInfoOverlay shows help text over the aquarium.
// It centers lines on screen so controls are easy to read.
func (a *Animation) drawInfoOverlay() {
	lines := InfoLines()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
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
			fg := infoStyleFor(i, ln, ch)
			termbox.SetCell(x+ci, y, ch, fg, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

// Run is the main state machine for the app.
// It initializes terminal mode, runs setup, then loops on input and ticks.
// Info mode pauses movement so overlay text stays easy to read.
func (a *Animation) Run(setup func(*Animation, bool), classic bool) error {
	if err := termbox.Init(); err != nil {
		return err
	}
	defer func() {
		termbox.Interrupt()
		termbox.Close()
	}()

	if err := a.updateSize(); err != nil {
		return err
	}
	a.running = true
	setup(a, classic)

	events := make(chan termbox.Event, 64)
	go func() {
		for {
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventInterrupt {
				return
			}
			events <- ev
		}
	}()

	paused := false
	showingInfo := false

	for a.running {
	drainEvents:
		for {
			select {
			case ev := <-events:
				switch ev.Type {
				case termbox.EventResize:
					if err := a.updateSize(); err != nil {
						return err
					}
					a.reflowForResize()
					if showingInfo {
						a.drawInfoOverlay()
					} else {
						a.drawFrame()
					}
				case termbox.EventKey:
					if ev.Key == termbox.KeyEsc && showingInfo {
						showingInfo = false
						paused = false
					} else if ev.Ch != 0 {
						switch ev.Ch {
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
				case termbox.EventError:
					if ev.Err != nil {
						return ev.Err
					}
				}
			default:
				break drainEvents
			}
		}
		if showingInfo {
			a.drawInfoOverlay()
		} else if !paused {
			a.animate()
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// EnsureScreenSupport checks if terminal drawing is available.
// It returns a clear error when termbox cannot initialize the terminal.
func EnsureScreenSupport() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	termbox.Close()
	return nil
}
