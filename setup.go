package main

// SetupAquarium builds the scene in a deliberate order.
// Environment/decor goes first, then long-living populations, then one event.
// This function is reused at startup and after reset.
func SetupAquarium(anim *Animation, classicMode bool) {
	AddEnvironment(anim)
	AddCastle(anim)
	AddAllSeaweed(anim)
	AddAllFish(anim, classicMode)
	RandomObject(nil, anim)
}
