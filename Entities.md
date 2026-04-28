# Visual Entity Catalog

## Overview

This document catalogs all visual entities rendered by `goquarium`, grouped into classes and described with implementation-level details.

Primary spawn flow:
- Initial scene setup happens in `SetupAquarium` (`setup.go`): environment, castle, seaweed, fish, then one random special object.
- Ongoing special entities are chained through `RandomObject` (`special.go`) via death callbacks.

An item is treated as a visual entity if it is created through `NewEntity`/`AddEntity` and rendered by the animation loop (`animation.go`), including composite helper entities such as hit points and linked parts.

---

## Class 1: Environment & Scenery

### `water_seg_0` to `water_seg_3`
- **Class:** Environment & Scenery
- **Defined/Spawned In:** `AddEnvironment` in `environment.go`
- **EntityType / Name:** `EntityType: "waterline"`, `Name: "water_seg_<i>"`
- **Visual Form:** Single-frame tiled water bands (`~`/`^` patterns) stretched across screen width.
- **Coloring:** `DefaultColor: "CYAN"`, no explicit color mask.
- **Motion:** Static (no movement callback).
- **Depth Layer:** `Depth["water_line0"]..Depth["water_line3"]` (fallback to `water_line0`).
- **Lifecycle:** Persistent scene entities.
- **Interactions:** `Physical: true`; bubbles collide with `waterline` and pop.

### `castle`
- **Class:** Environment & Scenery
- **Defined/Spawned In:** `AddCastle` in `environment.go`
- **EntityType / Name:** empty `EntityType`, `Name: "castle"`
- **Visual Form:** Large static multiline ASCII castle.
- **Coloring:** Explicit multiline mask (`R`, `W`, `y`, `w` markers) with `DefaultColor: "BLACK"`.
- **Motion:** Static.
- **Depth Layer:** `Depth["castle"]`.
- **Lifecycle:** Persistent scene entity.
- **Interactions:** Decorative only.

### `seaweed_<rand>`
- **Class:** Environment & Scenery
- **Defined/Spawned In:** `AddSeaweed` and `AddAllSeaweed` in `environment.go`
- **EntityType / Name:** empty `EntityType`, `Name: "seaweed_<random>"`
- **Visual Form:** Two-frame alternating plant shape built from `(` and `)` with randomized height (3-6).
- **Coloring:** `DefaultColor: "GREEN"`, no explicit mask.
- **Motion:** Frame animation only using `CallbackArgs [0, 0, 0, frameStep]`.
- **Depth Layer:** `Depth["seaweed"]`.
- **Lifecycle:** Timed (`DieTime`) and self-regenerating via `DeathCallback: AddSeaweed`.
- **Interactions:** Decorative only.

---

## Class 2: Core Marine Life

### `fish`
- **Class:** Core Marine Life
- **Defined/Spawned In:** `AddFish` and `AddAllFish` in `fish.go`
- **EntityType / Name:** `EntityType: "fish"`
- **Visual Form:** Single-frame per instance; shape selected from directional variants in `oldFishDesigns` or `newFishDesigns`.
- **Coloring:** Directional color mask from selected design, randomized by `randColor`.
- **Motion:** `Callback: FishCallback`; primary movement via `CallbackArgs [dx, 0, 0]`.
- **Depth Layer:** Random from `Depth["fish_start"]` to `Depth["fish_end"]`.
- **Lifecycle:** `DieOffscreen: true`; respawns through fish death callback to keep population stable.
- **Interactions:** `Physical: true`; collides with shark teeth and hook system (`hook_point`).

### Fish Visual Subtypes (`oldFishDesigns`, `newFishDesigns`)
- **Class:** Core Marine Life (subtype definitions for `fish`)
- **Defined/Spawned In:** design arrays in `fish.go`, selected in `AddFish`
- **EntityType / Name:** not separate runtime entity types
- **Visual Form:** directional sprite/mask pairs:
  - `oldFishDesigns`: 8 variants
  - `newFishDesigns`: 4 variants
- **Coloring:** numeric placeholders (`1..9`) replaced with random color markers (`c/C/r/R/y/Y/b/B/g/G/m/M`).
- **Motion:** Inherit `fish` behavior.
- **Depth Layer:** Inherit `fish` behavior.
- **Lifecycle:** Inherit `fish` behavior.
- **Interactions:** Inherit `fish` behavior.

### `bubble`
- **Class:** Core Marine Life
- **Defined/Spawned In:** `AddBubble` in `fish.go` (emitted from `FishCallback`)
- **EntityType / Name:** `EntityType: "bubble"`
- **Visual Form:** 5-frame animation: `.`, `o`, `O`, `O`, `O`.
- **Coloring:** `DefaultColor: "CYAN"`.
- **Motion:** `CallbackArgs [0, -1, 0, 0.1]` (rises upward with frame progression).
- **Depth Layer:** Spawned at fish depth minus one (`z-1`).
- **Lifecycle:** `DieOffscreen: true`.
- **Interactions:** `Physical: true`; collision handler kills bubble when touching `waterline`.

### `splat` (bite effect)
- **Class:** Core Marine Life (FX)
- **Defined/Spawned In:** `AddSplat` in `fish.go` (from `FishCollision`)
- **EntityType / Name:** empty `EntityType`
- **Visual Form:** 4-frame multiline burst effect.
- **Coloring:** `DefaultColor: "RED"`.
- **Motion:** No positional movement, frame stepping only (`CallbackArgs [0,0,0,0.25]`).
- **Depth Layer:** Around bite position, offset to `z-2`.
- **Lifecycle:** `DieFrame: 15` (short-lived).
- **Interactions:** Spawned when small fish are hit by `teeth`.

---

## Class 3: Special Event Creatures

### `ship`
- **Class:** Special Event Creatures
- **Defined/Spawned In:** `AddShip` in `special.go`
- **EntityType / Name:** empty `EntityType`
- **Visual Form:** Single-frame (directional) sailboat sprite.
- **Coloring:** Directional color mask for sails/hull, `DefaultColor: "WHITE"`.
- **Motion:** Horizontal drift via `CallbackArgs [±1, 0, 0, 0]`.
- **Depth Layer:** `Depth["water_gap1"]`.
- **Lifecycle:** `DieOffscreen: true`, chains into `RandomObject`.
- **Interactions:** Decorative moving special.

### `whale`
- **Class:** Special Event Creatures
- **Defined/Spawned In:** `AddWhale` in `special.go`
- **EntityType / Name:** empty `EntityType`
- **Visual Form:** 12 frames (5 idle + 7 spout frames), directional.
- **Coloring:** Directional mask reused across frames, `DefaultColor: "WHITE"`.
- **Motion:** `CallbackArgs [±0.5, 0, 0, 1]`.
- **Depth Layer:** `Depth["water_gap2"]`.
- **Lifecycle:** `DieOffscreen: true`, chains into `RandomObject`.
- **Interactions:** Decorative moving special.

### `monster` (new/old variants)
- **Class:** Special Event Creatures
- **Defined/Spawned In:** `AddMonster`, `addNewMonster`, `addOldMonster` in `special.go`
- **EntityType / Name:** empty `EntityType`
- **Visual Form:** Two monster families:
  - new monster: 2-frame animation
  - old monster: 4-frame animation
- **Coloring:** Eye highlight masks repeated per frame, `DefaultColor: "GREEN"`.
- **Motion:** `CallbackArgs [±2, 0, 0, 0.25]`.
- **Depth Layer:** `Depth["water_gap2"]`.
- **Lifecycle:** `DieOffscreen: true`, chains into `RandomObject`.
- **Interactions:** Event creature in random-object rotation.

### `big fish` (design1/design2 variants)
- **Class:** Special Event Creatures
- **Defined/Spawned In:** `AddBigFish`, `addBigFish1`, `addBigFish2` in `special.go`
- **EntityType / Name:** empty `EntityType`
- **Visual Form:** Large single-frame special fish with directional versions.
- **Coloring:** Large masks passed through `randColor`, `DefaultColor: "YELLOW"`.
- **Motion:** Horizontal movement:
  - design1: `CallbackArgs [±3.0, 0, 0]`
  - design2: `CallbackArgs [±2.5, 0, 0]`
- **Depth Layer:** `Depth["shark"]`.
- **Lifecycle:** `DieOffscreen: true`, chains into `RandomObject`.
- **Interactions:** Event creature in random-object rotation.

### `swan`
- **Class:** Special Event Creatures
- **Defined/Spawned In:** `AddSwan` in `special.go`
- **EntityType / Name:** empty `EntityType`
- **Visual Form:** Single-frame directional sprite.
- **Coloring:** Directional accent masks, `DefaultColor: "WHITE"`.
- **Motion:** `CallbackArgs [±1, 0, 0, 0.25]`.
- **Depth Layer:** `Depth["water_gap3"]`.
- **Lifecycle:** `DieOffscreen: true`, chains into `RandomObject`.
- **Interactions:** Decorative moving special.

### `ducks`
- **Class:** Special Event Creatures
- **Defined/Spawned In:** `AddDucks` in `special.go`
- **EntityType / Name:** empty `EntityType`
- **Visual Form:** 3-frame directional animated flock sprite.
- **Coloring:** Directional masks, `DefaultColor: "WHITE"`.
- **Motion:** `CallbackArgs [±1, 0, 0, 0.25]`.
- **Depth Layer:** `Depth["water_gap3"]`.
- **Lifecycle:** `DieOffscreen: true`, chains into `RandomObject`.
- **Interactions:** Decorative moving special.

### `dolphins` (formation)
- **Class:** Special Event Creatures
- **Defined/Spawned In:** `AddDolphins` in `special.go`
- **EntityType / Name:** empty `EntityType`
- **Visual Form:** Three separate entities in fixed spacing, each 2-frame directional animation.
- **Coloring:** Shared directional mask; defaults are `BLUE`, `BLUE`, and `CYAN` by formation index.
- **Motion:** `CallbackArgs [±2, 0, 0, 0.5]`.
- **Depth Layer:** `Depth["water_gap3"]`.
- **Lifecycle:** `DieOffscreen: true`; only lead dolphin has death callback to `RandomObject`.
- **Interactions:** Composite visual formation.

---

## Class 4: Interaction & Composite Systems

### Shark Composite System (`shark` + `teeth`)

#### `shark`
- **Class:** Interaction & Composite Systems
- **Defined/Spawned In:** `AddShark` in `special.go`
- **EntityType / Name:** `EntityType: "shark"`
- **Visual Form:** Large directional single-frame predator sprite.
- **Coloring:** Directional mask with `DefaultColor: "CYAN"`.
- **Motion:** `CallbackArgs [±2, 0, 0]`.
- **Depth Layer:** `Depth["shark"]`.
- **Lifecycle:** `DieOffscreen: true`; `DeathCallback: SharkDeath`.
- **Interactions:** Linked to `teeth` entity for collision-driven predation.

#### `teeth`
- **Class:** Interaction & Composite Systems
- **Defined/Spawned In:** `AddShark` in `special.go`
- **EntityType / Name:** `EntityType: "teeth"`
- **Visual Form:** Single `*` hitbox marker.
- **Coloring:** Default renderer color fallback (no explicit mask/default color override).
- **Motion:** `CallbackArgs [±2, 0, 0]`, matched with shark speed.
- **Depth Layer:** `Depth["shark"] + 1`.
- **Lifecycle:** Removed by `SharkDeath`.
- **Interactions:** `fish` collision handler uses `teeth` hits to trigger fish death and `splat`.

### Fishhook Composite System (`fishline` + `fishhook` + `hook_point`)

#### `fishline`
- **Class:** Interaction & Composite Systems
- **Defined/Spawned In:** `AddFishhook` in `special.go`
- **EntityType / Name:** `EntityType: "fishline"`
- **Visual Form:** Tall multiline line (`|\n` repeated with trailing blank tail).
- **Coloring:** Default renderer color fallback.
- **Motion:** `Callback: FishhookCallback` with state map args (`lowering`/`hooked`).
- **Depth Layer:** `Depth["water_line1"]`.
- **Lifecycle:** Removed by grouped death cleanup.
- **Interactions:** Retracted when fish is hooked.

#### `fishhook`
- **Class:** Interaction & Composite Systems
- **Defined/Spawned In:** `AddFishhook` in `special.go`
- **EntityType / Name:** `EntityType: "fishhook"`
- **Visual Form:** Multi-line hook sprite.
- **Coloring:** `DefaultColor: "GREEN"`.
- **Motion:** `Callback: FishhookCallback` with state map args.
- **Depth Layer:** `Depth["water_line1"]`.
- **Lifecycle:** `DieOffscreen: true`; death cleanup removes `hook_point` and `fishline`.
- **Interactions:** Part of hook system that can catch fish.

#### `hook_point`
- **Class:** Interaction & Composite Systems
- **Defined/Spawned In:** `AddFishhook` in `special.go`
- **EntityType / Name:** `EntityType: "hook_point"`
- **Visual Form:** Tiny 4-line marker used as collision point.
- **Coloring:** `DefaultColor: "GREEN"`.
- **Motion:** `Callback: FishhookCallback` with state map args.
- **Depth Layer:** `Depth["shark"] + 1`.
- **Lifecycle:** Removed in grouped hook cleanup.
- **Interactions:** `Physical: true`; fish collision with `hook_point` initiates retract sequence.

---

## Appendix A: Depth & Layer Reference

From `depth.go`:
- `gui_text: 0`, `gui: 1`
- `shark: 2`
- `fish_start: 3` to `fish_end: 20`
- `seaweed: 21`, `castle: 22`
- water bands and gaps:
  - `water_line3: 2`, `water_gap3: 3`
  - `water_line2: 4`, `water_gap2: 5`
  - `water_line1: 6`, `water_gap1: 7`
  - `water_line0: 8`, `water_gap0: 9`

Practical note: renderer sorts by depth and draws in descending order (`animation.go`), so layer relationships are controlled by these depth values.

## Appendix B: Shared Rendering & Animation Rules

- **Frame selection:** `Entity.CurrentFrame` indexes `Shapes`/`Colors` cyclically (`entity.go`).
- **Movement callback args convention:** `[]float64{dx, dy, dz, frameStep}` when default movement is used.
- **Auto transparency:** with `AutoTrans: true`, spaces/transparent runes are skipped in draw path (`animation.go`).
- **Color masks:** `CurrentColor` lines map per-character mask markers to colors through `maskColorMap` (`animation.go`).
- **Entity lifecycle controls:** `DieOffscreen`, `DieFrame`, `DieTime`, explicit `Kill`, and `DeathCallback`.

## Appendix C: Completeness Cross-Check

Cataloged visual entities instantiated from:
- `environment.go`: `waterline` segments, castle, seaweed
- `fish.go`: fish, bubble, splat
- `special.go`: ship, whale, monster variants, big fish variants, shark system, fishhook system, swan, ducks, dolphins

Supporting behavior/layer references verified against:
- `setup.go`, `depth.go`, `entity.go`, `animation.go`
