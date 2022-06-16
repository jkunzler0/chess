package database

import (
	"errors"
	"sync"
)

type Score struct {
	Win  int
	Loss int
}

var store = struct {
	sync.RWMutex
	m map[string]Score
}{m: make(map[string]Score)}

var ErrorNoSuchUser = errors.New("no such user")

func Put(user string, s Score) error {
	store.Lock()
	defer store.Unlock()

	store.m[user] = s

	return nil
}

func Delete(user string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.m, user)

	return nil
}

func IncrWinLoss(winner string, losser string) error {
	store.Lock()
	defer store.Unlock()

	value, ok := store.m[winner]
	if !ok {
		value = Score{Win: 0, Loss: 0}
	}
	value.Win += 1
	store.m[winner] = value

	value, ok = store.m[losser]
	if !ok {
		value = Score{Win: 0, Loss: 0}
	}
	value.Loss += 1
	store.m[losser] = value

	return nil
}

func Get(user string) (Score, error) {
	store.RLock()
	defer store.RUnlock()

	value, ok := store.m[user]

	if !ok {
		return value, ErrorNoSuchUser
	}

	return value, nil
}
