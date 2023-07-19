package control

import (
	"fmt"
	"net/http"
)

func (c *Control) Index(w http.ResponseWriter, r *http.Request) {
	if err := c.Renderer.Render(templateName("index", r), w, c.Accounts.Summary()); err != nil {
		fmt.Fprintf(w, "Could not render index: %v", err)
	}
}

func (c *Control) Save(w http.ResponseWriter, r *http.Request) {
	c.Accounts.Save(c.AccountsPath)
	if err := c.Renderer.Render(templateName("index", r), w, c.Accounts.Summary()); err != nil {
		fmt.Fprintf(w, "Could not render index: %v", err)
	}
}
