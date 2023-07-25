package control

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jwiklund/ah/history"
)

func (c *Control) Edit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Printf("Could not parse form: %v", err)
		}
		date := c.Accounts.CurrentDate()
		if strings.HasSuffix(r.URL.Path, "/amount") {
			key, value, err := slugAndIntValue(r)
			if err != nil {
				fmt.Printf("Could not parse amount: %v", err)
			} else {
				err = c.Accounts.UpdateAmountBySlug(key, date, value)
				if err != nil {
					fmt.Printf("Could not update amount: %v", err)
				}
			}
		}
		if strings.HasSuffix(r.URL.Path, "/change") {
			key, value, err := slugAndIntValue(r)
			if err != nil {
				fmt.Printf("Could not parse amount: %v", err)
			} else {
				err = c.Accounts.UpdateChangeBySlug(key, date, value)
				if err != nil {
					fmt.Printf("Could not update amount: %v", err)
				}
			}
		}
		if strings.HasSuffix(r.URL.Path, "/add") {
			name := r.Form.Get("name")
			if name != "" {
				c.Accounts.AddAccount(name, date)
			}
		}
	}
	if err := c.Renderer.Render(templateName("edit", r), w, c.Accounts.Current()); err != nil {
		fmt.Fprintf(w, "Couild not render edit: %v", err)
	}
}

type EditAccount struct {
	Name    string
	Slug    string
	History []history.SummaryEntry
	Total   Total
	Message string
	Error   error
}

func (c *Control) EditAccount(w http.ResponseWriter, r *http.Request) {
	edit := EditAccount{}
	render := func() {
		for _, h := range edit.History {
			edit.Total.Assets = h.End
			edit.Total.Change = edit.Total.Change + h.Change
			edit.Total.Increase = edit.Total.Increase + h.Increase
		}
		if err := c.Renderer.Render(templateName("edit.account", r), w, edit); err != nil {
			fmt.Fprintf(w, "Could not render edit account: %v", err)
		}
	}
	refreshRender := func() {
		edit.Name, edit.History, _ = c.Accounts.AccountHistory(edit.Slug)
		render()
	}

	path := strings.TrimPrefix(r.URL.Path, "/edit/account/")
	end := strings.Index(path, "/")
	if end == -1 {
		edit.Slug = path
	} else {
		edit.Slug = path[0:end]
		path = path[end+1:]
	}

	if edit.Slug == "" {
		edit.Error = errors.New("No account")
		render()
		return
	}

	if r.Method != "POST" {
		edit.Name, edit.History, edit.Error = c.Accounts.AccountHistory(edit.Slug)
		render()
		return
	}

	end = strings.Index(path, "/")
	if end == -1 {
		edit.Error = errors.New("No year")
		refreshRender()
	}

	year := path[0:end]
	path = path[end:]

	edit.Error = r.ParseForm()
	if edit.Error != nil {
		refreshRender()
		return
	}

	if path == "/amount" {
		amount, err := formIntInput(r, "amount")
		if err != nil {
			edit.Error = err
			refreshRender()
			return

		}
		edit.Error = c.Accounts.UpdateAmountBySlug(edit.Slug, year, amount)
	} else if path == "/change" {
		change, err := formIntInput(r, "change")
		if err != nil {
			edit.Error = err
			refreshRender()
			return
		}
		edit.Error = c.Accounts.UpdateChangeBySlug(edit.Slug, year, change)
	} else {
		edit.Error = fmt.Errorf("Unknown update for path %s", path)
	}
	refreshRender()
}
