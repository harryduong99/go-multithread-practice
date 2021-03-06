package deadlock

import (
	"time"

	"github.com/duongnam99/go-multithread-practice/deadlocksTrain/common"
)

func MoveTrain(train *common.Train, distance int, crossings []*common.Crossing) {
	for train.Front < distance {
		train.Front += 1
		// check wether we are at the crossing or not
		for _, crossing := range crossings {
			if train.Front == crossing.Position {
				crossing.Intersection.Mutex.Lock()
				crossing.Intersection.LockedBy = train.Id
			}
			back := train.Front - train.TrainLength
			if back == crossing.Position { // if train passed the crossing
				crossing.Intersection.LockedBy = -1
				crossing.Intersection.Mutex.Unlock()
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}
