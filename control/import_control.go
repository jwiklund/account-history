package control

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func (c *Control) Import(w http.ResponseWriter, r *http.Request) {
	historyData := HistoryData{
		Separators:  csv.ColumnSeparators,
		ColumnTypes: csv.ColumnTypes,
		Options: csv.ImportOptions{
			Separator: csv.SpaceLike,
		},
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Printf("Could not parse form: %v", err)
		}
		separator := csv.ColumnSeparator(rawFormInput(r, "separator"))
		if _, ok := csv.ColumnSeparators[separator]; ok {
			historyData.Options.Separator = separator
		}
		historyData.Csv = formInput(r, "csv")
		historyData.Options.Name = formInput(r, "name")
		historyData.Options.Date = formInput(r, "date")
		if !strings.HasSuffix(r.URL.Path, "/separator") {
			historyData.Options.Columns = importPostColumns(r)
		}
		csvData := formInput(r, "csv")
		if csvData != "" {
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
			if err == nil {
				historyData.Error = history.Update(historyData.Options, c.Accounts)
				if historyData.Error == nil {
					historyData.Message = fmt.Sprintf("Updated %d rows", len(history.Rows))
				}
			}
		}
		if isHx(r) {
			if err := c.Renderer.Render(templateNameWithPart("import", r, "import"), w, historyData); err != nil {
				fmt.Fprintf(w, "Could not render import: %v", err)
			}
			return
		}
	}
	if err := c.Renderer.Render(templateName("import", r), w, historyData); err != nil {
		fmt.Fprintf(w, "Could not render import: %v", err)
	}
}

func importPostColumns(r *http.Request) []csv.ImportColumnType {
	var result []csv.ImportColumnType
	duplicates := make(map[csv.ImportColumnType]int)
	i := 0
	for {
		currentName := strings.Join([]string{"column", strconv.Itoa(i)}, "-")
		current := formInput(r, currentName)
		if current == "" {
			return result
		}
		columnType := csv.ImportColumnType(current)
		if _, ok := csv.ColumnTypes[columnType]; ok {
			if j, ok := duplicates[columnType]; ok {
				if strings.HasSuffix(r.URL.Path, currentName) {
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
