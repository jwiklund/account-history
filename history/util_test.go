package history

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortAccount(t *testing.T) {
	accounts := []Account{
		makeAccount("b", "2022", 1, false),
		makeAccount("a", "2022", 1, false),
	}
	assert.Equal(t, []Account{
		makeAccount("a", "2022", 1, false),
		makeAccount("b", "2022", 1, false),
	}, sortAccounts(accounts))
}

func TestSortAccountOneoff(t *testing.T) {
	accounts := []Account{
		makeAccount("a", "2022", 1, true),
		makeAccount("b", "2022", 1, false),
	}
	assert.Equal(t, []Account{
		makeAccount("b", "2022", 1, false),
		makeAccount("a", "2022", 1, true),
	}, sortAccounts(accounts))
}

func TestSortAccountEmpty(t *testing.T) {
	accounts := []Account{
		makeAccount("a", "2022", 0, false),
		makeAccount("b", "2022", 1, false),
	}
	assert.Equal(t, []Account{
		makeAccount("b", "2022", 1, false),
		makeAccount("a", "2022", 0, false),
	}, sortAccounts(accounts))
}

func makeAccount(name, year string, end int, oneoff bool) Account {
	return Account{
		Name:   name,
		Oneoff: oneoff,
		History: []History{
			{
				Date:   year,
				Amount: end,
			},
		},
	}
}
