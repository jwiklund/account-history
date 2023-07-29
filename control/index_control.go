package control

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/jwiklund/ah/history"
)

func (c *Control) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tag := r.URL.Query().Get("tag")
	if err := c.Renderer.Render(templateName("index", r), w, summarize(*c.Accounts, tag)); err != nil {
		fmt.Fprintf(w, "Could not render index: %v", err)
	}
}

func (c *Control) Save(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c.Accounts.Save(c.AccountsPath)
	tag := r.URL.Query().Get("tag")
	if err := c.Renderer.Render(templateName("index", r), w, summarize(*c.Accounts, tag)); err != nil {
		fmt.Fprintf(w, "Could not render index: %v", err)
	}
}

type IndexData struct {
	Years []history.SummaryEntry
	Total Total
	Tag   string
	Tags  []string
}

type Total struct {
	Assets   int
	Increase int
	Change   int
}

func summarize(a history.Accounts, tag string) IndexData {
	data := IndexData{
		Tag: tag,
	}
	data.Years, data.Tags = a.Summary(data.Tag)

	var totalSum int
	var totalIncrease int
	var totalChange int

	for _, y := range data.Years {
		totalSum = y.End
		totalIncrease = totalIncrease + y.Increase
		totalChange = totalChange + y.Change
	}

	data.Total = Total{
		Assets:   totalSum,
		Increase: totalIncrease,
		Change:   totalChange,
	}
	return data
}
