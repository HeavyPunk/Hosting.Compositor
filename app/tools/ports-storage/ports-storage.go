package ports_storage

import (
	"errors"
	"log"
	"os/exec"
	"time"
)

// Работает только с TCP (в UDP пока нет смысла)
// Временно будет сделано с кешем в памяти, периодически обновляемым через netstat
// Потом переедет в sqlite
// Потоконебезопасно!!!

var Free = true
var Busy = false

var portCache = make(map[int]bool)
var updatePeriod = 3 * time.Minute
var lastUpdateTimestamp = time.Now()

var minPort = 10001
var maxPort = 10100

func OccupyPort(port int) error {
	if time.Now().Sub(lastUpdateTimestamp) > updatePeriod {
		portCache = Update()
		lastUpdateTimestamp = time.Now()
	}
	return nil
}

func ReleasePort(port int) error {
	if checkForHardwarePort(port) {
		portCache[port] = Free
		return nil
	} else {
		return errors.New("Port cannot be released due handling by some process")
	}
}

func GetRandomFreePort() (int, error) {
	for i := minPort; i < maxPort; i++ {
		if IsFreePort(i) {
			return i, nil
		}
	}
	return -1, errors.New("All available ports are busy")
}

func IsFreePort(port int) bool {
	isFree, ok := portCache[port]
	return isFree || !ok
}

func checkForHardwarePort(port int) bool {
	cmd := exec.Command("netstat", "-t", "-l")
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	busyPorts := extractPortsFromNetstatOut(stdout)
	for _, el := range busyPorts {
		if el == port {
			return Busy
		}
	}
	return Free
}

func extractPortsFromNetstatOut(out []byte) []int {
	arr := make([]int, 3)
	arr[0] = 80
	arr[1] = 443
	arr[2] = 8080
	return arr
}

func Update() map[int]bool {
	cmd := exec.Command("netstat", "-t", "-l")
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	busyPorts := extractPortsFromNetstatOut(stdout)
	cache := make(map[int]bool)
	for _, el := range busyPorts {
		cache[el] = Busy
	}
	return cache
}
