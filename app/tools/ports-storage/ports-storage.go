package ports_storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"simple-hosting/compositor/app/settings"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/semaphore"
)

// Работает только с TCP (в UDP пока нет смысла)
// Временно будет сделано с кешем в памяти, периодически обновляемым через netstat
// Потом переедет в sqlite
// Потоконебезопасно!!!

type PortsStorageContext struct {
	portCache  map[int]bool
	config     settings.ServiceSettings
	dbDriver   string
	dbFilePath string
	minPort    int
	maxPort    int
	sem        *semaphore.Weighted
	ctx        context.Context
}

var (
	Free = true
	Busy = false
)

func Init(config settings.ServiceSettings) *PortsStorageContext {
	return &PortsStorageContext{
		portCache:  make(map[int]bool),
		config:     config,
		dbDriver:   config.Hypervisor.Services.PortsService.DbDriver,
		dbFilePath: config.Hypervisor.Services.PortsService.DbPath,
		minPort:    config.Hypervisor.Services.PortsService.MinPort,
		maxPort:    config.Hypervisor.Services.PortsService.MaxPort,
		sem:        semaphore.NewWeighted(1),
		ctx:        context.TODO(),
	}
}

func (context *PortsStorageContext) occupyPort(port int) error {
	isFree, err := context.isFreePort(port)
	if err != nil {
		return err
	}
	if !isFree {
		return errors.New("port " + strconv.Itoa(port) + "is already occupied")
	}
	context.portCache[port] = Busy
	db, err := sql.Open(context.dbDriver, context.dbFilePath)
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

func (context *PortsStorageContext) ReleasePort(port int) error {
	context.portCache[port] = Free
	db, err := sql.Open(context.dbDriver, context.dbFilePath)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("delete from ports where port=$1", port)
	return err
}

func (context *PortsStorageContext) GetRandomFreePort() (int, error) {
	if err := context.sem.Acquire(context.ctx, 1); err != nil {
		return -1, errors.New("acquire semaphore failed")
	}
	defer context.sem.Release(1)

	for i := context.minPort; i < context.maxPort; i++ {
		isFree, err := context.isFreePort(i)
		if err != nil {
			return i, err
		}
		if isFree {
			err := context.occupyPort(i)
			return i, err
		}
	}
	return -1, errors.New("all available ports are busy")
}

func (context *PortsStorageContext) isFreePort(port int) (bool, error) {
	isFree, ok := context.portCache[port]
	portFreeInCache := isFree || !ok
	if portFreeInCache {
		portIsFree, err := context.checkInDatabase(port)
		if err != nil {
			return false, err
		}
		context.portCache[port] = portIsFree
		return portIsFree, nil
	}
	return portFreeInCache, nil
}

func (context *PortsStorageContext) checkInDatabase(port int) (bool, error) {
	db, err := sql.Open(context.dbDriver, context.dbFilePath)
	if err != nil {
		return false, err
	}
	defer db.Close()
	res, err := db.Query("select port from ports where port=$1", port)
	if err != nil {
		return false, err
	}
	defer res.Close()
	for res.Next() {
		p := 0
		if err = res.Scan(&p); err != nil {
			fmt.Println(err)
			continue
		}
		if p == port {
			return Busy, nil
		}
	}
	return Free, nil
}
