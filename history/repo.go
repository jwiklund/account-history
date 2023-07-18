package history

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

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

func Load(filename string, initHistory bool) (*Accounts, error) {
	reader, err := os.Open(filename)
	if os.IsNotExist(err) && initHistory {
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

func (a *Accounts) Save(filename string) error {
	writer, err := os.Create(filename)
	if err != nil {
		return err
	}
	encoder := yaml.NewEncoder(writer)
	for _, account := range a.accounts {
		err := encoder.Encode(account)
		if err != nil {
			return fmt.Errorf("could not write entry to %s: %w", filename, err)
		}
	}
	return encoder.Close()
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
	Slug     string
	Start    int
	End      int
	Change   int
	Increase int
}

func (a *Accounts) Current() []CurrentEntry {
	date := a.CurrentDate()
	var current []CurrentEntry
	for _, a := range a.accounts {
		lastIndex := len(a.History) - 1
		if a.History[lastIndex].Date != date {
			current = append(current, CurrentEntry{
				Name:     a.Name,
				Slug:     nameToSlug(a.Name),
				Start:    a.History[lastIndex].Amount,
				End:      a.History[lastIndex].Amount,
				Change:   0,
				Increase: 0,
			})
		} else if lastIndex == 0 {
			current = append(current, CurrentEntry{
				Name:     a.Name,
				Slug:     nameToSlug(a.Name),
				Start:    0,
				End:      a.History[lastIndex].Amount,
				Change:   a.History[lastIndex].Change,
				Increase: a.History[lastIndex].Amount - a.History[lastIndex].Change,
			})
		} else {
			current = append(current, CurrentEntry{
				Name:     a.Name,
				Slug:     nameToSlug(a.Name),
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

func (a *Accounts) CurrentDate() string {
	var date string
	for _, a := range a.accounts {
		for _, h := range a.History {
			if date == "" || strings.Compare(date, h.Date) < 0 {
				date = h.Date
			}
		}
	}
	if date != "" {
		return date
	}
	return time.Now().Format("2006")
}

func (a *Accounts) AddAccount(name string, date string) {
	a.accounts = append(a.accounts, Account{
		Name: name,
		History: []History{
			{Date: date},
		},
	})
}

func (a *Accounts) UpdateAmountBySlug(slug string, date string, newAmount int) error {
	return a.updateBySlug(slug, date, func(h History) History {
		h.Amount = newAmount
		return h
	})
}

func (a *Accounts) UpdateChangeBySlug(slug string, date string, newChange int) error {
	return a.updateBySlug(slug, date, func(h History) History {
		h.Change = newChange
		return h
	})
}

func (a *Accounts) updateBySlug(slug string, date string, update func(History) History) error {
	for _, account := range a.accounts {
		if nameToSlug(account.Name) == slug {
			for index, history := range account.History {
				if history.Date == date {
					account.History[index] = update(history)
					return nil
				}
			}
			account.History = append(account.History, update(History{
				Date:   date,
				Amount: 0,
				Change: 0,
			}))
			return nil
		}
	}
	return fmt.Errorf("No such account: %s", slug)
}

var nameToSlugRegex = regexp.MustCompile("[^a-z0-9]+")

func nameToSlug(name string) string {
	return nameToSlugRegex.ReplaceAllString(strings.ToLower(name), "-")
}
