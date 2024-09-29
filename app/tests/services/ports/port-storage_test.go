package port_storage_test

import (
	"fmt"
	ports_storage "simple-hosting/compositor/app/tools/ports-storage"
	"testing"

	ports_test_base "simple-hosting/compositor/app/tests/services/ports/base"

	_ "github.com/mattn/go-sqlite3"
)

func TestOccupyPort(t *testing.T) {
	settings := ports_test_base.GetSettings("settings.yml")
	err := ports_test_base.PrepareDatabase(settings)
	if err != nil {
		t.Error(err)
	}

  defer func() {
    err := ports_test_base.DisposeDatabase(settings)
    if err != nil {
      t.Error(err)
    }
  }()

	storage := ports_storage.Init(settings)
	channels := make(map[int]chan error)
	opCount := 5
	for i := 0; i < opCount; i++ {
		channels[i] = make(chan error)
	}

	for i := 0; i < opCount; i++ {
		go func(ch chan error) {
			port, err := storage.GetRandomFreePort()
			if err != nil {
				ch <- err
				return
			}

			if port <= 0 {
				ch <- fmt.Errorf("returned port %d is out of range", port)
				return
			}
			ch <- nil
		}(channels[i])
	}

	for i := 0; i < opCount; i++ {
		if err = <-channels[i]; err != nil {
			t.Errorf("Error from thread %d: %v", i, err)
		}
	}
}
