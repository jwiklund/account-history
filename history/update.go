package history

import (
	"fmt"
)

func (a *Accounts) AddAccount(name string, date string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	slug := NameToSlug(name)
	for _, a := range a.accounts {
		if NameToSlug(a.Name) == slug {
			return fmt.Errorf("Account %s collides with %s", name, a.Name)
		}
	}

	a.accounts = append(a.accounts, Account{
		Name: name,
		History: []History{
			{Date: date},
		},
	})
	sortAccounts(a.accounts)
	return nil
}

func (a *Accounts) AddEmptyAccount(name string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	slug := NameToSlug(name)
	for _, a := range a.accounts {
		if NameToSlug(a.Name) == slug {
			return fmt.Errorf("Account %s collides with %s", name, a.Name)
		}
	}

	a.accounts = append(a.accounts, Account{Name: name})
	sortAccounts(a.accounts)
	return nil
}

func (a *Accounts) UpdateAmountBySlug(slug string, date string, newAmount int) error {
	return a.UpdateHistoryBySlugDate(slug, date, func(h History) History {
		h.Amount = newAmount
		return h
	})
}

func (a *Accounts) UpdateChangeBySlug(slug string, date string, newChange int) error {
	return a.UpdateHistoryBySlugDate(slug, date, func(h History) History {
		h.Change = newChange
		return h
	})
}

func (a *Accounts) UpdateHistoryBySlugDate(slug string, date string, update func(History) History) error {
	return a.UpdateHistoryBySlug(slug, func(history []History) ([]History, error) {
		for index, entry := range history {
			if entry.Date == date {
				history[index] = update(entry)
				return history, nil
			}
		}
		return append(history, update(History{
			Date:   date,
			Amount: 0,
			Change: 0,
		})), nil
	})
}

func (a *Accounts) UpdateHistoryBySlug(slug string, update func([]History) ([]History, error)) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	for i, account := range a.accounts {
		if NameToSlug(account.Name) == slug {
			newHistory, err := update(account.History)
			if err == nil {
				sortHistory(newHistory)
				account.History = newHistory
				a.accounts[i] = account
			}
			return err
		}
	}
	return fmt.Errorf("No such account: %s", slug)
}
