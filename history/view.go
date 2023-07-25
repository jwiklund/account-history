package history

import (
	"fmt"
	"sort"
)

func (a *Accounts) AccountHistory(slug string) (string, []SummaryEntry, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for _, a := range a.accounts {
		if NameToSlug(a.Name) != slug {
			continue
		}
		var summary []SummaryEntry
		current := 0
		for _, h := range a.History {
			summary = append(summary, SummaryEntry{
				Year:     h.Date,
				Start:    current,
				End:      h.Amount,
				Change:   h.Change,
				Increase: h.Amount - current - h.Change,
			})
			current = h.Amount
		}
		return a.Name, summary, nil
	}
	return "", nil, fmt.Errorf("No such account: %s", slug)
}

type SummaryEntry struct {
	Year     string
	Start    int
	End      int
	Change   int
	Oneoff   int
	Increase int
}

func (a *Accounts) Summary() []SummaryEntry {
	a.lock.Lock()
	defer a.lock.Unlock()

	summary := make(map[string]*SummaryEntry)
	for _, a := range a.accounts {
		for _, h := range a.History {
			entry, ok := summary[h.Date]
			if !ok {
				entry = &SummaryEntry{
					Year: h.Date,
				}
				summary[h.Date] = entry
			}
			entry.End = entry.End + h.Amount
			if a.Oneoff {
				entry.Oneoff = entry.Oneoff - h.Change
				entry.Change = entry.Change + h.Change
			} else {
				entry.Change = entry.Change + h.Change
			}
		}
	}
	var dates []string
	for date := range summary {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	var result []SummaryEntry
	current := 0
	for _, date := range dates {
		entry := summary[date]
		entry.Start = current
		entry.Increase = entry.End - entry.Change - entry.Start - entry.Oneoff
		current = entry.End
		result = append(result, *entry)
	}
	return result
}

type CurrentEntry struct {
	Name     string
	Slug     string
	Start    int
	End      int
	Change   int
	Increase int
}

func (a *Accounts) Current() []CurrentEntry {
	a.lock.Lock()
	defer a.lock.Unlock()

	date := currentDateLocked(a.accounts)
	var current []CurrentEntry
	for _, a := range a.accounts {
		lastIndex := len(a.History) - 1
		if lastIndex == -1 {
			current = append(current, CurrentEntry{
				Name:     a.Name,
				Slug:     NameToSlug(a.Name),
				Start:    0,
				End:      0,
				Change:   0,
				Increase: 0,
			})
		} else if a.History[lastIndex].Date != date {
			current = append(current, CurrentEntry{
				Name:     a.Name,
				Slug:     NameToSlug(a.Name),
				Start:    a.History[lastIndex].Amount,
				End:      a.History[lastIndex].Amount,
				Change:   0,
				Increase: 0,
			})
		} else if lastIndex == 0 {
			current = append(current, CurrentEntry{
				Name:     a.Name,
				Slug:     NameToSlug(a.Name),
				Start:    0,
				End:      a.History[lastIndex].Amount,
				Change:   a.History[lastIndex].Change,
				Increase: a.History[lastIndex].Amount - a.History[lastIndex].Change,
			})
		} else {
			current = append(current, CurrentEntry{
				Name:     a.Name,
				Slug:     NameToSlug(a.Name),
				Start:    a.History[lastIndex-1].Amount,
				End:      a.History[lastIndex].Amount,
				Change:   a.History[lastIndex].Change,
				Increase: a.History[lastIndex].Amount - a.History[lastIndex].Change - a.History[lastIndex-1].Amount,
			})
		}
	}
	return current
}

func (a *Accounts) CurrentDate() string {
	a.lock.Lock()
	defer a.lock.Unlock()

	return currentDateLocked(a.accounts)
}
