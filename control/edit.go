package control

import (
	"fmt"
	"net/http"
	"strings"
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
