package control

import (
	"fmt"
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/jwiklund/ah/history"
)

func (c *Control) Index(w http.ResponseWriter, r *http.Request) {
	if err := c.Renderer.Render(templateName("index", r), w, summarize(*c.Accounts)); err != nil {
		fmt.Fprintf(w, "Could not render index: %v", err)
	}
}

func (c *Control) Save(w http.ResponseWriter, r *http.Request) {
	c.Accounts.Save(c.AccountsPath)
	if err := c.Renderer.Render(templateName("index", r), w, summarize(*c.Accounts)); err != nil {
		fmt.Fprintf(w, "Could not render index: %v", err)
	}
}

type IndexData struct {
	Years []history.SummaryEntry
	Total []NameValue
}

type NameValue struct {
	Name  string
	Value string
}

func summarize(a history.Accounts) IndexData {
	var data IndexData
	data.Years = a.Summary()

	var totalSum float64
	var totalIncrease float64
	var totalChange float64

	for _, y := range data.Years {
		totalSum = float64(y.End)
		totalIncrease = totalIncrease + float64(y.Increase)
		totalChange = totalChange + float64(y.Change)
	}

	data.Total = []NameValue{
		{
			Name:  "Total assets",
			Value: humanize.SIWithDigits(totalSum, 2, " kr"),
		},
		{
			Name:  "Total increase",
			Value: humanize.SIWithDigits(totalIncrease, 2, " kr"),
		},
		{
			Name:  "Total change",
			Value: humanize.SIWithDigits(totalChange, 2, " kr"),
		},
	}
	return data
}
