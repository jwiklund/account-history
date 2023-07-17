package history

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Account struct {
	Name    string    `yaml:"name"`
	History []History `yaml:"history"`
}

type History struct {
	Date   string
	Amount int
	Change int
}

type Accounts struct {
	accounts []Account
}

func Load(filename string, initAccounts bool) (*Accounts, error) {
	reader, err := os.Open(filename)
	if err == os.ErrNotExist && initAccounts {
		return &Accounts{nil}, nil
	}
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	result, err := LoadFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", filename, err)
	}
	return result, err
}

func LoadFrom(reader io.Reader) (*Accounts, error) {
	decoder := yaml.NewDecoder(reader)
	var result []Account
	for {
		var account Account
		err := decoder.Decode(&account)
		if err == io.EOF {
			return &Accounts{result}, nil
		}
		if err != nil {
			return nil, err
		}
		result = append(result, account)
	}
}

type SummaryEntry struct {
	Year     string
	Start    int
	End      int
	Change   int
	Increase int
}

func (a *Accounts) Summary() []SummaryEntry {
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
			entry.Change = entry.Change + h.Change
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
		entry.Increase = entry.End - entry.Change - entry.Start
		current = entry.End
		result = append(result, *entry)
	}
	return result
}

type CurrentEntry struct {
	Name     string
	Start    int
	End      int
	Change   int
	Increase int
}

func (a *Accounts) Current() []CurrentEntry {
	var date string
	for _, a := range a.accounts {
		for _, h := range a.History {
			if date == "" || strings.Compare(date, h.Date) < 0 {
				date = h.Date
			}
		}
	}

	var current []CurrentEntry
	for _, a := range a.accounts {
		lastIndex := len(a.History) - 1
		if a.History[lastIndex].Date != date {
			current = append(current, CurrentEntry{
				Name:     a.Name,
				Start:    a.History[lastIndex].Amount,
				End:      a.History[lastIndex].Amount,
				Change:   0,
				Increase: 0,
			})
		} else if lastIndex == 0 {
			current = append(current, CurrentEntry{
				Name:     a.Name,
				Start:    0,
				End:      a.History[lastIndex].Amount,
				Change:   a.History[lastIndex].Change,
				Increase: a.History[lastIndex].Amount - a.History[lastIndex].Change,
			})
		} else {
			current = append(current, CurrentEntry{
				Name:     a.Name,
				Start:    a.History[lastIndex-1].Amount,
				End:      a.History[lastIndex].Amount,
				Change:   a.History[lastIndex].Change,
				Increase: a.History[lastIndex].Amount - a.History[lastIndex].Change - a.History[lastIndex-1].Amount,
			})
		}
	}
	sort.Slice(current, func(i, j int) bool {
		firstEmpty := current[i].End == 0 && current[i].Change == 0
		secondEmpty := current[j].End == 0 && current[j].Change == 0
		bothEmpty := firstEmpty && secondEmpty
		neitherEmpty := !firstEmpty && !secondEmpty
		if bothEmpty || neitherEmpty {
			return current[i].Name < current[j].Name
		}
		return secondEmpty
	})
	return current
}
