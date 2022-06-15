package database

import (
	"errors"
	"sync"
)

type Score struct {
	win  int
	loss int
}

var store = struct {
	sync.RWMutex
	m map[string]Score
}{m: make(map[string]Score)}

var ErrorNoSuchKey = errors.New("no such key")

func Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.m, key)

	return nil
}

func Get(key string) (Score, error) {
	store.RLock()
	defer store.RUnlock()

	value, ok := store.m[key]

	if !ok {
		return value, ErrorNoSuchKey
	}

	return value, nil
}

func Put(key string, s Score) error {
	store.Lock()
	defer store.Unlock()

	store.m[key] = s

	return nil
}

func IncrementWinLoss(winner string, losser string) error {
	store.Lock()
	defer store.Unlock()

	value, ok := store.m[winner]
	if !ok {
		value = Score{win: 0, loss: 0}
	}
	value.win += 1
	store.m[winner] = value

	value, ok = store.m[losser]
	if !ok {
		value = Score{win: 0, loss: 0}
	}
	value.loss += 1
	store.m[losser] = value

	return nil
}
