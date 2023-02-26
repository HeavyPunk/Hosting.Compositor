package ports_storage

import (
	"errors"
	"time"
)

// Работает только с TCP (в UDP пока нет смысла)
// Временно будет сделано с кешем в памяти, периодически обновляемым через netstat
// Потом переедет в sqlite

var portCache = make(map[int]bool)
var updatePeriod = 3 * time.Minute
var lastUpdateTimestamp = time.Now()

func OccupyPort(port int) error {
	if time.Now().Sub(lastUpdateTimestamp) > updatePeriod {
		Update()
		lastUpdateTimestamp = time.Now()
	}
	return nil
}

func ReleasePort(port int) error {
	if checkForHardwarePort(port) {
		portCache[port] = true
		return nil
	} else {
		return errors.New("Port cannot be released due handling by some process")
	}
}

func IsFreePort(port int) bool {
	return portCache[port]
}

func checkForHardwarePort(port int) bool {
	panic("Not implemented")
}

func Update() map[int]bool {
	panic("Not implemented")
}
