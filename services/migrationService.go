package services

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type MigrationService struct {
	dataProvider WorkerDataProvider
}

func (ms *MigrationService) initMigrationService(provider WorkerDataProvider) *MigrationService {
	ms.dataProvider = provider
	return ms
}

func NewMigrationService(provider WorkerDataProvider) *MigrationService {
	return (&MigrationService{}).initMigrationService(provider)
}

func (ms *MigrationService) updateCurrentHour() error {
	return ms.dataProvider.UpdateHour(time.Now().Truncate(time.Hour).Add(1 * time.Hour).Unix())
}

func (ms *MigrationService) Start() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	// update expenses for current hour every minute
	go func() {
		defer wg.Done()
		for {
			ms.updateCurrentHour()
			time.Sleep(time.Minute)
		}
	}()
	wg.Add(1)
	// go back in time and calculate expenses for the past
	go func() {
		defer wg.Done()

		// limit of how far in past we want it to go
		lastHourToUpdateString := os.Getenv("LAST_HOUR_TO_UPDATE")
		lastHourToUpdate, err := strconv.Atoi(lastHourToUpdateString)
		if err != nil {
			fmt.Println(err)
			return
		}

		furthestUpdatedHour := time.Now().Truncate(time.Hour).Unix()
		for {
			hour, err := ms.dataProvider.GetFurthestUpdatedHour()
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Minute)
			} else {
				furthestUpdatedHour = hour
				break
			}
		}
		// we already processed all we wanted
		if furthestUpdatedHour <= int64(lastHourToUpdate) {
			return
		}
		for {
			hourToUpdate := furthestUpdatedHour - 3600
			err := ms.dataProvider.UpdateHour(hourToUpdate)
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Minute)
			}
			furthestUpdatedHour = hourToUpdate
		}
	}()
	wg.Wait()
}
