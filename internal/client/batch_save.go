/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package client

import (
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type SaveTask struct {
	PlayerID int
	SaveFunc func() error
}

type BatchSaveManager struct {
	tasks       map[int]*SaveTask
	mutex       sync.RWMutex
	batchSize   int
	batchTicker *time.Ticker
}

var globalBatchSaveManager *BatchSaveManager

func init() {
	globalBatchSaveManager = NewBatchSaveManager(50, 30*time.Second) // Process 50 saves every 30 seconds
	go globalBatchSaveManager.Start()
}

func NewBatchSaveManager(batchSize int, interval time.Duration) *BatchSaveManager {
	return &BatchSaveManager{
		tasks:       make(map[int]*SaveTask),
		batchSize:   batchSize,
		batchTicker: time.NewTicker(interval),
	}
}

func (bsm *BatchSaveManager) QueueSave(playerID int, saveFunc func() error) {
	bsm.mutex.Lock()
	defer bsm.mutex.Unlock()

	// Overwrite existing task if any (avoid duplicate saves)
	bsm.tasks[playerID] = &SaveTask{
		PlayerID: playerID,
		SaveFunc: saveFunc,
	}
}

func (bsm *BatchSaveManager) Start() {
	for {
		select {
		case <-bsm.batchTicker.C:
			bsm.processBatch()
		}
	}
}

func (bsm *BatchSaveManager) processBatch() {
	bsm.mutex.Lock()

	if len(bsm.tasks) == 0 {
		bsm.mutex.Unlock()
		return
	}

	// Take a batch of tasks
	count := 0
	tasksToProcess := make([]*SaveTask, 0, bsm.batchSize)
	for _, task := range bsm.tasks {
		tasksToProcess = append(tasksToProcess, task)
		delete(bsm.tasks, task.PlayerID)
		count++
		if count >= bsm.batchSize {
			break
		}
	}

	bsm.mutex.Unlock()

	// Process tasks concurrently with goroutine limit
	semaphore := make(chan struct{}, 10) // Max 10 concurrent saves
	var wg sync.WaitGroup

	for _, task := range tasksToProcess {
		wg.Add(1)
		go func(t *SaveTask) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := t.SaveFunc(); err != nil {
				log.Errorf("Batch save failed for player %d: %v", t.PlayerID, err)
			}
		}(task)
	}

	wg.Wait()
	log.Debugf("Processed batch of %d saves", len(tasksToProcess))
}

// Global function to access the batch save manager
func QueueBatchSave(playerID int, saveFunc func() error) {
	globalBatchSaveManager.QueueSave(playerID, saveFunc)
}
