package main

import (
	"fmt"
	"runtime"
)

// InfoText builds the long help screen text.
// It is shown when user runs --info.
// This gives a friendly overview of controls and purpose.
func InfoText() string {
	return fmt.Sprintf(`
╔═══════════════════════════════════════════════════════════════════════╗
║                                                                       ║
║   Asciiquarium %s - ASCII Art Aquarium Animation                      ║
║                                                                       ║
╚═══════════════════════════════════════════════════════════════════════╝

An aquarium/sea animation in ASCII art for your terminal!

CONTROLS:
  Q or q  - Quit the aquarium
  P or p  - Pause/unpause animation
  R or r  - Redraw and respawn entities
  I or i  - Show/hide info screen (press I or ESC to return)
`, Version)
}

// InfoLines returns a short line-by-line overlay message.
// The animation screen draws these lines in the center.
// We keep this separate so rendering code stays simple.
func InfoLines() []string {
	return []string{
		"╔═══════════════════════════════════════════════════════════════════════╗",
		"║                                                                       ║",
		fmt.Sprintf("║   Asciiquarium %s - ASCII Art Aquarium Animation                      ║", Version),
		"║                                                                       ║",
		"╚═══════════════════════════════════════════════════════════════════════╝",
		"",
		"  Q/q quit   P/p pause   R/r reset   I/i info   ESC close info",
		"",
		"  Press I or ESC to return to aquarium...",
	}
}

// VersionString builds one compact version line.
// It includes app version, Go runtime, and OS/arch.
// This is printed for --version and -v.
func VersionString() string {
	return fmt.Sprintf("goquarium/%s Go/%s %s/%s", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
