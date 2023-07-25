package history

import (
	"regexp"
	"sort"
	"strings"
	"time"
)

var nameToSlugRegex = regexp.MustCompile("[^a-z0-9]+")

func NameToSlug(name string) string {
	return nameToSlugRegex.ReplaceAllString(strings.ToLower(name), "-")
}

func currentDateLocked(accounts []Account) string {
	var date string
	for _, a := range accounts {
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

func sortHistory(history []History) {
	sort.Slice(history, func(i, j int) bool {
		return history[i].Date < history[j].Date
	})
}

func sortAccounts(accounts []Account) []Account {
	histories := make([]History, 0, len(accounts))
	for _, a := range accounts {
		if len(a.History) == 0 {
			histories = append(histories, History{})
		} else {
			histories = append(histories, a.History[len(a.History)-1])
		}
	}
	indexes := make([]int, len(accounts))
	for i := range accounts {
		indexes[i] = i
	}
	sort.Slice(indexes, func(i, j int) bool {
		firstOneoff := accounts[indexes[i]].Oneoff
		secondOneoff := accounts[indexes[j]].Oneoff
		bothOneoff := firstOneoff && secondOneoff
		eitherOneoff := firstOneoff || secondOneoff
		if !bothOneoff && eitherOneoff {
			return secondOneoff
		}

		firstEmpty := histories[indexes[i]].Amount == 0 && histories[indexes[i]].Change == 0
		secondEmpty := histories[indexes[j]].Amount == 0 && histories[indexes[j]].Change == 0
		bothEmpty := firstEmpty && secondEmpty
		eitherEmpty := firstEmpty || secondEmpty
		if !bothEmpty && eitherEmpty {
			return secondEmpty
		}

		return accounts[indexes[i]].Name < accounts[indexes[j]].Name
	})
	result := make([]Account, 0, len(accounts))
	for _, i := range indexes {
		result = append(result, accounts[i])
	}
	return result
}
