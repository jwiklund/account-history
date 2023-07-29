package csv

import (
	"testing"

	"github.com/jwiklund/ah/history"
	"github.com/stretchr/testify/assert"
)

var amount1Rows = ImportRows{
	[]ImportRow{
		{
			Columns: []ImportColumn{
				{
					Value: "1",
					Type:  Amount,
				},
			},
		},
	},
}

var opts = ImportOptions{
	Name: "name",
	Date: "2022",
}

func TestUpdateEmpty(t *testing.T) {
	empty := history.New()
	err := amount1Rows.Update(opts, empty)
	assert.Error(t, err, "No such account: name")
}

func TestUpdateEmptyWithAccount(t *testing.T) {
	empty := history.New()
	empty.AddEmptyAccount("name")
	err := amount1Rows.Update(opts, empty)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	summary, _ := empty.Summary("")
	assert.Equal(t, []history.SummaryEntry{{
		Year:     "2022",
		Start:    0,
		End:      1,
		Change:   0,
		Increase: 1,
	}}, summary)
	assert.Equal(t, []history.CurrentEntry{{
		Name:     "name",
		Slug:     "name",
		Start:    0,
		End:      1,
		Change:   0,
		Increase: 1,
	}}, empty.Current())
}
