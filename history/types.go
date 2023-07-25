package history

import "sync"

type Accounts struct {
	accounts []Account
	lock     *sync.Mutex
}

type Account struct {
	Name    string    `yaml:"name"`
	History []History `yaml:"history"`
}

type History struct {
	Date   string
	Amount int
	Change int
}

func New() *Accounts {
	return &Accounts{
		lock: &sync.Mutex{},
	}
}
