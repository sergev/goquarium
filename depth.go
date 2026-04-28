package main

// Depth stores drawing layers used by moving objects.
// Bigger/smaller numbers decide what appears in front.
// This helps keep fish, water, and props ordered correctly.
var Depth = map[string]int{
	"gui_text":    0,
	"gui":         1,
	"shark":       2,
	"fish_start":  3,
	"fish_end":    20,
	"seaweed":     21,
	"castle":      22,
	"water_line3": 2,
	"water_gap3":  3,
	"water_line2": 4,
	"water_gap2":  5,
	"water_line1": 6,
	"water_gap1":  7,
	"water_line0": 8,
	"water_gap0":  9,
}
