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

func sortAccounts(accounts []Account) {
	current := currentDateLocked(accounts)
	histories := make([]History, 0, len(accounts))
	for _, a := range accounts {
		if len(a.History) == 0 {
			histories = append(histories, History{})
		} else if a.History[len(a.History)-1].Date < current {
			histories = append(histories, History{})
		} else {
			histories = append(histories, a.History[len(a.History)-1])
		}
	}
	sort.Slice(accounts, func(i, j int) bool {
		firstEmpty := histories[i].Amount == 0 && histories[i].Change == 0
		secondEmpty := histories[j].Amount == 0 && histories[j].Change == 0
		bothEmpty := firstEmpty && secondEmpty
		neitherEmpty := !firstEmpty && !secondEmpty
		if bothEmpty || neitherEmpty {
			return accounts[i].Name < accounts[j].Name
		}
		return secondEmpty
	})
}
