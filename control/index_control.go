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

	var totalSum int64
	var totalIncrease int64
	var totalChange int64

	for _, y := range data.Years {
		totalSum = int64(y.End)
		totalIncrease = totalIncrease + int64(y.Increase)
		totalChange = totalChange + int64(y.Change)
	}

	data.Total = []NameValue{
		{
			Name:  "Total assets",
			Value: humanize.Comma(totalSum),
		},
		{
			Name:  "Total increase",
			Value: humanize.Comma(totalIncrease),
		},
		{
			Name:  "Total change",
			Value: humanize.Comma(totalChange),
		},
	}
	return data
}
