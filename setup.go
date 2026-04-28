package main

// SetupAquarium creates the first scene of the animation.
// It adds environment objects, fish, and one random special object.
// This is called at startup and when user presses reset.
func SetupAquarium(anim *Animation, classicMode bool) {
	AddEnvironment(anim)
	AddCastle(anim)
	AddAllSeaweed(anim)
	AddAllFish(anim, classicMode)
	RandomObject(nil, anim)
}
