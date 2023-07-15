package history

import (
	"bytes"
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
	assert.Equal(t, &Accounts{[]Account{
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
	}}, accounts)
}

func TestSummary(t *testing.T) {
	accounts := &Accounts{
		[]Account{
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
	summary := accounts.Summary()
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
