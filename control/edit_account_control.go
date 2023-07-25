package control

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/jwiklund/ah/history"
)

type EditAccount struct {
	Name    string
	Slug    string
	History []history.SummaryEntry
	Total   Total
	Message string
	Error   error
}

func (c *Control) RenderEditAccount(w http.ResponseWriter, r *http.Request, slug, message string, err error) {
	edit := EditAccount{
		Slug:    slug,
		Message: message,
		Error:   err,
	}
	edit.Name, edit.History, err = c.Accounts.AccountHistory(slug)
	if edit.Error == nil {
		edit.Error = err
	}
	for _, h := range edit.History {
		edit.Total.Assets = h.End
		edit.Total.Change = edit.Total.Change + h.Change
		edit.Total.Increase = edit.Total.Increase + h.Increase
	}
	if err := c.Renderer.Render(templateName("edit.account", r), w, edit); err != nil {
		fmt.Fprintf(w, "Could not render edit account: %v", err)
	}
}

func (c *Control) EditAccount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	slug := p.ByName("accountSlug")
	c.RenderEditAccount(w, r, slug, "", nil)
}

func (c *Control) EditAccountAdd(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	slug := p.ByName("accountSlug")
	err := r.ParseForm()
	if err != nil {
		c.RenderEditAccount(w, r, slug, "", err)
		return
	}
	year := formInput(r, "year")
	err = c.Accounts.UpdateHistoryBySlugDate(slug, year, func(h history.History) history.History {
		return h
	})
	c.RenderEditAccount(w, r, slug, "", err)
}

func (c *Control) EditAccountAmount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	slug := p.ByName("accountSlug")
	year := p.ByName("year")
	err := r.ParseForm()
	if err != nil {
		c.RenderEditAccount(w, r, slug, "", err)
		return
	}
	amount, err := formIntInput(r, strings.Join([]string{year, "amount"}, "-"))
	if err != nil {
		c.RenderEditAccount(w, r, slug, "", err)
		return
	}
	err = c.Accounts.UpdateAmountBySlug(slug, year, amount)
	c.RenderEditAccount(w, r, slug, "", err)
}

func (c *Control) EditAccountChange(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	slug := p.ByName("accountSlug")
	year := p.ByName("year")
	err := r.ParseForm()
	if err != nil {
		c.RenderEditAccount(w, r, slug, "", err)
		return
	}
	change, err := formIntInput(r, strings.Join([]string{year, "change"}, "-"))
	if err != nil {
		c.RenderEditAccount(w, r, slug, "", err)
		return
	}
	err = c.Accounts.UpdateChangeBySlug(slug, year, change)
	c.RenderEditAccount(w, r, slug, "", err)
}
