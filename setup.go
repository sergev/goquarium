package main

func SetupAquarium(anim *Animation, classicMode bool) {
	AddEnvironment(anim)
	AddCastle(anim)
	AddAllSeaweed(anim)
	AddAllFish(anim, classicMode)
	RandomObject(nil, anim)
}
