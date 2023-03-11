package ports_storage

import (
	"database/sql"
	"errors"
	"fmt"
	"simple-hosting/compositor/app/settings"
	file_settings_provider "simple-hosting/go-commons/settings/file-settings-provider"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// Работает только с TCP (в UDP пока нет смысла)
// Временно будет сделано с кешем в памяти, периодически обновляемым через netstat
// Потом переедет в sqlite
// Потоконебезопасно!!!

var Free = true
var Busy = false

var portCache = make(map[int]bool)
var config = file_settings_provider.GetSetting[settings.ServiceSettings]("settings.yml")

var dbDriver = "sqlite3"
var dbFilePath = config.Hypervisor.Services.PortsService.DbPath

var minPort = 10001
var maxPort = 10100

func OccupyPort(port int) error {
	if !IsFreePort(port) {
		return errors.New("Port " + strconv.Itoa(port) + " already occupied")
	}
	portCache[port] = Busy
	db, err := sql.Open(dbDriver, dbFilePath)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("insert into ports (port) values ($1)", port)
	if err != nil {
		return err
	}
	return nil
}

func ReleasePort(port int) error {
	portCache[port] = Free
	db, err := sql.Open(dbDriver, dbFilePath)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("delete from ports where port=$1", port)
	return err
}

func GetRandomFreePort() (int, error) {
	for i := minPort; i < maxPort; i++ {
		if IsFreePort(i) {
			err := OccupyPort(i)
			return i, err
		}
	}
	return -1, errors.New("All available ports are busy")
}

func IsFreePort(port int) bool {
	isFree, ok := portCache[port]
	portFreeInCache := isFree || !ok
	if portFreeInCache {
		portIsFree := checkInDatabase(port)
		portCache[port] = portIsFree
		return portIsFree
	}
	return portFreeInCache
}

func checkInDatabase(port int) bool {
	db, err := sql.Open(dbDriver, dbFilePath)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query("select port from ports where port=$1", port)
	if err != nil {
		panic(err)
	}
	defer res.Close()
	for res.Next() {
		p := 0
		if err = res.Scan(&p); err != nil {
			fmt.Println(err)
			continue
		}
		if p == port {
			return Busy
		}
	}
	return Free
}

// func checkForHardwarePort(port int) bool {
// 	cmd := exec.Command("bash", config.Hypervisor.Services.ScriptsDir+"print-listen-ports.sh")
// 	stdout, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		log.Fatal(err)
// 	}
// 	busyPorts := extractPortsFromLsofOut(stdout)
// 	for _, el := range busyPorts {
// 		if el == port {
// 			return Busy
// 		}
// 	}
// 	return Free
// }

// func Update() map[int]bool {
// 	cmd := exec.Command("netstat", "-t", "-l")
// 	stdout, err := cmd.Output()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	busyPorts := extractPortsFromLsofOut(stdout)
// 	cache := make(map[int]bool)
// 	for _, el := range busyPorts {
// 		cache[el] = Busy
// 	}
// 	return cache
// }
