package control

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/jwiklund/ah/history"
)

type Edit struct {
	Current []history.CurrentEntry
	Message string
	Error   error
}

func (c *Control) RenderEdit(w http.ResponseWriter, r *http.Request, message string, err error) {
	edit := Edit{
		Current: c.Accounts.Current(),
		Message: message,
		Error:   err,
	}
	if err := c.Renderer.Render(templateName("edit", r), w, edit); err != nil {
		fmt.Fprintf(w, "Couild not render edit: %v", err)
	}
}

func (c *Control) Edit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c.RenderEdit(w, r, "", nil)
}

func (c *Control) EditAdd(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		c.RenderEdit(w, r, "", err)
		return
	}
	date := c.Accounts.CurrentDate()
	name := r.Form.Get("name")
	if name != "" {
		err = c.Accounts.AddAccount(name, date)
	}
	c.RenderEdit(w, r, "", err)
}

func (c *Control) EditAmount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		c.RenderEdit(w, r, "", err)
		return
	}

	date := c.Accounts.CurrentDate()
	key, value, err := slugAndIntValue(r)
	if err == nil {
		err = c.Accounts.UpdateAmountBySlug(key, date, value)
	}
	c.RenderEdit(w, r, "", err)
}

func (c *Control) EditChange(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		c.RenderEdit(w, r, "", err)
		return
	}

	date := c.Accounts.CurrentDate()
	key, value, err := slugAndIntValue(r)
	if err == nil {
		err = c.Accounts.UpdateChangeBySlug(key, date, value)
	}
	c.RenderEdit(w, r, "", err)
}
