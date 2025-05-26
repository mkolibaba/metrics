package jsonfile

import (
	"context"
	"fmt"
	"github.com/mkolibaba/metrics/internal/server/storage/inmemory"
	"go.uber.org/zap"
	"os"
	"time"
)

type FileDatabase interface {
	Save(gauges map[string]float64, counters map[string]int64) error
	Load() (map[string]float64, map[string]int64, error)
	Close()
}

type FileStorage struct {
	*inmemory.MemStorage
	db          FileDatabase
	instantSync bool
	logger      *zap.SugaredLogger
}

func (f *FileStorage) UpdateGauge(ctx context.Context, name string, value float64) (float64, error) {
	if f.instantSync {
		defer f.save()
	}
	return f.MemStorage.UpdateGauge(ctx, name, value)
}

func (f *FileStorage) UpdateCounter(ctx context.Context, name string, value int64) (int64, error) {
	if f.instantSync {
		defer f.save()
	}
	return f.MemStorage.UpdateCounter(ctx, name, value)
}

func NewFileStorage(path string, storeInterval time.Duration, shouldRestore bool, logger *zap.SugaredLogger) (*FileStorage, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", path, err)
	}

	db := newFileDB(file)
	delegateStore, err := newUnderlyingStorage(db, shouldRestore)
	if err != nil {
		return nil, fmt.Errorf("error creating store: %v", err)
	}

	store := &FileStorage{
		MemStorage: delegateStore,
		db:         db,
		logger:     logger,
	}

	if storeInterval > 0 {
		go func() {
			for {
				time.Sleep(storeInterval)
				store.save()
			}
		}()
	} else if storeInterval == 0 {
		store.instantSync = true
	} else {
		return nil, fmt.Errorf("error creating store: storeInterval must be non-negative")
	}

	return store, nil
}

func (f *FileStorage) save() {
	gauges, err := f.GetGauges(context.TODO())
	if err != nil {
		f.logger.Errorf("error retrieving gauges for saving: %v", err)
		return
	}
	counters, err := f.GetCounters(context.TODO())
	if err != nil {
		f.logger.Errorf("error retrieving counters for saving: %v", err)
		return
	}

	if err := f.db.Save(gauges, counters); err != nil {
		f.logger.Errorf("error saving metrics: %v", err)
	}
}

func (f *FileStorage) Close() {
	f.db.Close()
}

func restore(db FileDatabase, store *inmemory.MemStorage) error {
	gauges, counters, err := db.Load()
	if err != nil {
		return err
	}

	for k, v := range counters {
		_, err := store.UpdateCounter(context.TODO(), k, v)
		if err != nil {
			return err
		}
	}
	for k, v := range gauges {
		_, err := store.UpdateGauge(context.TODO(), k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func newUnderlyingStorage(db FileDatabase, shouldRestore bool) (*inmemory.MemStorage, error) {
	store := inmemory.NewMemStorage()
	if shouldRestore {
		if err := restore(db, store); err != nil {
			return nil, err
		}
	}
	return store, nil
}
