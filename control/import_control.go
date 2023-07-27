package control

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/jwiklund/ah/csv"
)

type HistoryData struct {
	Csv         string
	Separators  map[csv.ColumnSeparator]string
	ColumnTypes map[csv.ImportColumnType]string
	Options     csv.ImportOptions
	History     []csv.ImportRow
	Error       error
	Message     string
}

func (c *Control) RenderImport(w http.ResponseWriter, r *http.Request, historyData HistoryData) {
	if err := c.Renderer.Render(templateName("import", r), w, historyData); err != nil {
		fmt.Fprintf(w, "Could not render import: %v", err)
	}
}

func (c *Control) RenderPartialImport(w http.ResponseWriter, r *http.Request, historyData HistoryData) {
	if isHx(r) {
		if err := c.Renderer.Render(templateNameWithPart("import", r, "import"), w, historyData); err != nil {
			fmt.Fprintf(w, "Could not render import: %v", err)
		}
		return
	}
	c.RenderImport(w, r, historyData)
}

func (c *Control) Import(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	historyData := HistoryData{
		Separators:  csv.ColumnSeparators,
		ColumnTypes: csv.ColumnTypes,
		Options: csv.ImportOptions{
			Separator: csv.SpaceLike,
			Plugins:   c.ImportPlugins,
		},
	}
	c.RenderImport(w, r, historyData)
}

func (c *Control) PrepareImport(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	historyData := c.prepareImportData(r, "")
	if historyData.Error != nil {
		c.RenderPartialImport(w, r, historyData)
	}
	historyData = prepareImportCsv(r, historyData)
	c.RenderPartialImport(w, r, historyData)
}

func (c *Control) PrepareImportSeparator(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	historyData := c.prepareImportData(r, "")
	if historyData.Error != nil {
		c.RenderPartialImport(w, r, historyData)
	}
	historyData.Options.Columns = nil
	historyData = prepareImportCsv(r, historyData)
	c.RenderPartialImport(w, r, historyData)
}

func (c *Control) PrepareImportPlugin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	historyData := c.prepareImportData(r, "")
	if historyData.Error != nil {
		c.RenderPartialImport(w, r, historyData)
	}
	historyData.Options.Separator = csv.SpaceLike
	historyData.Options.Columns = nil
	historyData.Options.Name = ""
	historyData.Options.Date = ""
	historyData.Options.CurrentDate = c.Accounts.CurrentDate()
	historyData = prepareImportCsv(r, historyData)
	c.RenderPartialImport(w, r, historyData)
}

func (c *Control) PrepareImportColumn(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	historyData := c.prepareImportData(r, p.ByName("columnId"))
	if historyData.Error != nil {
		c.RenderPartialImport(w, r, historyData)
	}
	historyData = prepareImportCsv(r, historyData)
	c.RenderPartialImport(w, r, historyData)
}

func (c *Control) ImportData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	historyData := c.prepareImportData(r, "")
	if historyData.Error != nil {
		c.RenderPartialImport(w, r, historyData)
	}
	historyData = prepareImportCsv(r, historyData)
	if historyData.Error == nil && historyData.History != nil {
		historyData.Error = csv.ImportRows{Rows: historyData.History}.Update(historyData.Options, c.Accounts)
		if historyData.Error == nil {
			historyData.Message = fmt.Sprintf("Updated %d rows", len(historyData.History))
		}
	}
	c.RenderPartialImport(w, r, historyData)
}

func (c *Control) prepareImportData(r *http.Request, columnId string) HistoryData {
	historyData := HistoryData{
		Separators:  csv.ColumnSeparators,
		ColumnTypes: csv.ColumnTypes,
		Options: csv.ImportOptions{
			Separator: csv.SpaceLike,
			Plugins:   c.ImportPlugins,
		},
	}
	historyData.Error = r.ParseForm()
	if historyData.Error != nil {
		return historyData
	}
	separator := csv.ColumnSeparator(rawFormInput(r, "separator"))
	if _, ok := csv.ColumnSeparators[separator]; ok {
		historyData.Options.Separator = separator
	}
	plugin := formInput(r, "plugin")
	if _, ok := c.ImportPlugins[plugin]; ok {
		historyData.Options.Plugin = plugin
	}

	historyData.Csv = formInput(r, "csv")
	historyData.Options.Name = formInput(r, "name")
	historyData.Options.Date = formInput(r, "date")
	historyData.Options.Columns = importPostColumns(r, columnId)
	return historyData
}

func prepareImportCsv(r *http.Request, historyData HistoryData) HistoryData {
	csvData := formInput(r, "csv")
	if csvData == "" {
		return historyData
	}
	lines := strings.Split(csvData, "\n")
	history, columns, name, date, err := csv.Import(lines, historyData.Options)
	historyData.History = history.Rows
	historyData.Options.Columns = columns
	historyData.Error = err
	if name != "" {
		historyData.Options.Name = name
	}
	if date != "" {
		historyData.Options.Date = date
	}
	return historyData
}

func importPostColumns(r *http.Request, changed string) []csv.ImportColumnType {
	var result []csv.ImportColumnType
	duplicates := make(map[csv.ImportColumnType]int)
	i := 0
	for {
		currentString := strconv.Itoa(i)
		currentName := strings.Join([]string{"column", currentString}, "-")
		current := formInput(r, currentName)
		if current == "" {
			return result
		}
		columnType := csv.ImportColumnType(current)
		if _, ok := csv.ColumnTypes[columnType]; ok {
			if j, ok := duplicates[columnType]; ok {
				if currentString == changed {
					result[j] = csv.None
				} else {
					columnType = csv.None
				}
			}
			result = append(result, columnType)
			if columnType != csv.None {
				duplicates[columnType] = i
			}
		} else {
			result = append(result, csv.None)
		}
		i = i + 1
	}
}
