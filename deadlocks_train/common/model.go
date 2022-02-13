package common

import "sync"

type Train struct {
	Id          int
	TrainLength int
	Front       int // position of the train
}

type Intersection struct { // where to track are passing
	Id       int
	Mutex    sync.Mutex // garantee that only one train is using the intersection
	LockedBy int        // which train is using that particular intersection
}

type Crossing struct { // whre the intersection is located along the track the the train is following
	Position     int
	Intersection *Intersection
}
