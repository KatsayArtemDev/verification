package processing

import "time"

func TimeChecking(dbTime time.Time) time.Duration {
	currentTime := time.Now().UTC()
	return currentTime.Sub(dbTime)
}
