package csv

import (
	"strconv"

	"github.com/jwiklund/ah/history"
	"golang.org/x/exp/slices"
)

type historyUpdate struct {
	date      string
	amount    int
	hasAmount bool
	change    int
	hasChange bool
}

type historyUpdates []historyUpdate

func (h ImportRows) Update(opts ImportOptions, accounts *history.Accounts) error {
	for slug, updates := range h.rowsBySlug(opts) {
		slices.SortFunc(updates, func(a, b historyUpdate) bool {
			return a.date < b.date
		})
		err := accounts.UpdateHistoryBySlug(slug, updates.update)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h ImportRows) rowsBySlug(opts ImportOptions) map[string]historyUpdates {
	rowsBySlug := make(map[string]historyUpdates)
	for _, row := range h.Rows {
		slug := history.NameToSlug(opts.Name)
		update := historyUpdate{
			date: opts.Date,
		}

		for _, column := range row.Columns {
			switch column.Type {
			case Name:
				slug = history.NameToSlug(column.Value)
			case Date:
				update.date = column.Value
			case Amount:
				update.amount, _ = strconv.Atoi(column.Value)
				update.hasAmount = true
			case Change:
				update.change, _ = strconv.Atoi(column.Value)
				update.hasChange = true
			}
		}

		if rows, ok := rowsBySlug[slug]; ok {
			rowsBySlug[slug] = append(rows, update)
		} else {
			rowsBySlug[slug] = []historyUpdate{update}
		}
	}
	return rowsBySlug
}

func (u historyUpdates) update(h []history.History) ([]history.History, error) {
	slices.SortFunc(h, func(a, b history.History) bool {
		return a.Date < b.Date
	})
	historyIndex := 0
	var result []history.History
	for updateIndex := 0; updateIndex < len(u); {
		if historyIndex < len(h) && h[historyIndex].Date < u[updateIndex].date {
			result = append(result, h[historyIndex])
			historyIndex++
			continue
		}
		if historyIndex < len(h) && h[historyIndex].Date == u[updateIndex].date {
			result = append(result, updateHistory(h[historyIndex], u[updateIndex]))
			historyIndex++
			updateIndex++
			continue
		}
		if historyIndex >= len(h) || h[historyIndex].Date > u[updateIndex].date {
			result = append(result, updateHistory(
				history.History{
					Date: u[updateIndex].date,
				},
				u[updateIndex]),
			)
			updateIndex++
		}
	}
	return result, nil
}

func updateHistory(h history.History, update historyUpdate) history.History {
	if update.hasAmount {
		h.Amount = update.amount
	}
	if update.hasChange {
		h.Change = update.change
	}
	return h
}
