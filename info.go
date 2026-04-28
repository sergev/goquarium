package main

import (
	"fmt"
	"runtime"
)

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

func VersionString() string {
	return fmt.Sprintf("goquarium/%s Go/%s %s/%s", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
