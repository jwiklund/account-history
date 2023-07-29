package history

import (
	"bytes"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const loadFromExample = `
name: name-1
history:
- date: 2022
  amount: 1
  change: 1
---
name: name-2
history:
- date: 2022
  amount: 1
  change: 1
- date: 2023
  amount: 2
  change: 0
`

func TestLoadFrom(t *testing.T) {
	accounts, err := LoadFrom(bytes.NewBuffer([]byte(loadFromExample)))
	assert.NoError(t, err)
	assert.Equal(t, []Account{
		{
			Name: "name-1",
			History: []History{
				{
					Date:   "2022",
					Amount: 1,
					Change: 1,
				},
			},
		},
		{
			Name: "name-2",
			History: []History{
				{
					Date:   "2022",
					Amount: 1,
					Change: 1,
				},
				{
					Date:   "2023",
					Amount: 2,
					Change: 0,
				},
			},
		},
	}, accounts.accounts)
}

func TestSummary(t *testing.T) {
	accounts := &Accounts{
		lock: &sync.Mutex{},
		accounts: []Account{
			{
				Name: "name-1",
				History: []History{
					{
						Date:   "2022",
						Amount: 1,
						Change: 1,
					},
					{
						Date:   "2023",
						Amount: 1,
						Change: 0,
					},
				},
			},
			{
				Name: "name-2",
				History: []History{
					{
						Date:   "2022",
						Amount: 1,
						Change: 1,
					},
					{
						Date:   "2023",
						Amount: 2,
						Change: 0,
					},
				},
			},
		},
	}
	summary, _ := accounts.Summary("")
	assert.Equal(t, []SummaryEntry{
		{
			Year:     "2022",
			Start:    0,
			End:      2,
			Change:   2,
			Increase: 0,
		},
		{
			Year:     "2023",
			Start:    2,
			End:      3,
			Change:   0,
			Increase: 1,
		},
	}, summary)
}

func TestCurrent(t *testing.T) {
	unsortedAccounts := []Account{
		{
			Name: "name-0",
			History: []History{
				{
					Date:   "2022",
					Amount: 0,
					Change: 0,
				},
			},
		},
		{
			Name: "name-0-2",
			History: []History{
				{
					Date:   "2022",
					Amount: 1,
					Change: 1,
				},
				{
					Date:   "2023",
					Amount: 0,
					Change: -1,
				},
			},
		},
		{
			Name: "name-1",
			History: []History{
				{
					Date:   "2022",
					Amount: 1,
					Change: 1,
				},
			},
		},
		{
			Name: "name-2",
			History: []History{
				{
					Date:   "2023",
					Amount: 2,
					Change: 0,
				},
			},
		},
		{
			Name: "name 3",
			History: []History{
				{
					Date:   "2022",
					Amount: 2,
					Change: 2,
				},
				{
					Date:   "2023",
					Amount: 2,
					Change: 0,
				},
			},
		},
	}
	accounts := &Accounts{
		lock:     &sync.Mutex{},
		accounts: sortAccounts(unsortedAccounts),
	}
	summary := accounts.Current()
	assert.Equal(t, []CurrentEntry{
		{
			Name:     "name 3",
			Slug:     "name-3",
			Start:    2,
			End:      2,
			Change:   0,
			Increase: 0,
		},
		{
			Name:     "name-0-2",
			Slug:     "name-0-2",
			Start:    1,
			End:      0,
			Change:   -1,
			Increase: 0,
		},
		{
			Name:     "name-1",
			Slug:     "name-1",
			Start:    1,
			End:      1,
			Change:   0,
			Increase: 0,
		},
		{
			Name:     "name-2",
			Slug:     "name-2",
			Start:    0,
			End:      2,
			Change:   0,
			Increase: 2,
		},
		{
			Name:     "name-0",
			Slug:     "name-0",
			Start:    0,
			End:      0,
			Change:   0,
			Increase: 0,
		},
	}, summary)
}
