package history

import "sync"

const (
	Oneoff string = "Oneoff"
)

type Accounts struct {
	accounts []Account
	lock     *sync.Mutex
}

type Account struct {
	Name    string    `yaml:"name"`
	History []History `yaml:"history"`
	Tags    []string  `yaml:"tags"`
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
