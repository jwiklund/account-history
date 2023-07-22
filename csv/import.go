package csv

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type ColumnSeparator string

const (
	SpaceLike ColumnSeparator = "[ \t]+"
	Space                     = " "
	Tab                       = "\t"
	Colon                     = ","
	SemiColon                 = ";"
)

var ColumnSeparators = map[ColumnSeparator]string{
	SpaceLike: "SpaceLike [ \\t]+",
	Space:     "Space [ ]",
	Tab:       "Tab [\\t]",
	Colon:     "Colon [:]",
	SemiColon: "SemiColon [;]",
}

type ImportColumnType string

const (
	None   ImportColumnType = ""
	Name   ImportColumnType = "name"
	Date   ImportColumnType = "date"
	Amount ImportColumnType = "amount"
	Change ImportColumnType = "change"
)

var ColumnTypes = map[ImportColumnType]string{
	None:   "None",
	Name:   "Name",
	Date:   "Date",
	Amount: "Amount",
	Change: "Change",
}

type ImportOptions struct {
	Separator ColumnSeparator
	Columns   []ImportColumnType
	Name      string
	Date      string
}

type ImportColumn struct {
	Value string
	Type  ImportColumnType
	Error error
}

type importColumnFeature struct {
	Numeric bool
	Date    bool
}

type ImportRow struct {
	Columns []ImportColumn
}

type ImportRows struct {
	Rows []ImportRow
}

func Import(csv []string, opts ImportOptions) (ImportRows, []ImportColumnType, string, string, error) {
	separator, err := regexp.Compile(string(opts.Separator))
	if err != nil {
		return ImportRows{}, nil, "", "", fmt.Errorf("invalid separator: %w", err)
	}
	var lines [][]ImportColumn
	var columnFeatures []importColumnFeature
	for _, line := range csv {
		if line == "" {
			continue
		}
		var columns []ImportColumn
		columns, columnFeatures = importLine(line, separator, columnFeatures)
		lines = append(lines, columns)
	}
	rows, columnTypes, name, date := typeRows(opts, lines, columnFeatures)
	return ImportRows{rows}, columnTypes, name, date, validateRows(rows, opts.Name, opts.Date)
}

func importLine(line string, separator *regexp.Regexp, columnFeatures []importColumnFeature) ([]ImportColumn, []importColumnFeature) {
	var columns []ImportColumn
	for i, value := range separator.Split(line, -1) {
		value = strings.TrimSpace(value)
		columns = append(columns, ImportColumn{
			Value: value,
		})
		_, err := strconv.Atoi(value)
		columnNumeric := err == nil
		columnDate := err == nil && len(value) == 4

		if len(columnFeatures) <= i {
			columnFeatures = append(columnFeatures, importColumnFeature{
				Numeric: columnNumeric,
				Date:    columnDate,
			})
		} else {
			if !columnNumeric {
				columnFeatures[i].Numeric = false
			}
			if !columnDate {
				columnFeatures[i].Date = false
			}
		}
	}
	return columns, columnFeatures
}

func typeRows(opts ImportOptions, lines [][]ImportColumn, features []importColumnFeature) ([]ImportRow, []ImportColumnType, string, string) {
	columnTypes := typeColumns(opts, features)
	name := opts.Name
	sameName := true
	date := opts.Date
	sameDate := true
	dynamicNameIndex := slices.Index(columnTypes, Name)
	dynamicDateIndex := slices.Index(columnTypes, Date)
	var rows []ImportRow
	for i := range lines {
		for j := range lines[i] {
			lines[i][j].Type = columnTypes[j]
			lines[i][j].Error = lines[i][j].Type.Validate(lines[i][j].Value)
		}
		if dynamicNameIndex != -1 {
			var currentName string
			if dynamicNameIndex < len(lines[i]) {
				currentName = lines[i][dynamicNameIndex].Value
			}
			if currentName == "" {

			} else if name == "" {
				name = currentName
			} else {
				if name != currentName {
					sameName = false
				}
			}
		}
		if dynamicDateIndex != -1 {
			var currentDate string
			if dynamicDateIndex < len(lines[i]) {
				currentDate = lines[i][dynamicDateIndex].Value
			}
			if currentDate == "" {
				sameDate = false
			} else if date == "" {
				date = currentDate
			} else {
				if date != lines[i][dynamicDateIndex].Value {
					sameDate = false
				}
			}
		}
		rows = append(rows, ImportRow{lines[i]})
	}
	if !sameName {
		name = opts.Name
	}
	if !sameDate {
		date = opts.Date
	}
	return rows, columnTypes, name, date
}

func typeColumns(opts ImportOptions, features []importColumnFeature) []ImportColumnType {
	hasName := opts.Name != "" || slices.Contains(opts.Columns, Name)
	hasDate := opts.Date != "" || slices.Contains(opts.Columns, Date)
	hasAmount := slices.Contains(opts.Columns, Amount)
	hasChange := slices.Contains(opts.Columns, Change)
	var result []ImportColumnType
	for i, feature := range features {
		if i < len(opts.Columns) && opts.Columns[i] != None {
			result = append(result, opts.Columns[i])
			continue
		}
		if feature.Date && !hasDate {
			hasDate = true
			result = append(result, Date)
			continue
		}
		if feature.Numeric {
			if !hasAmount {
				hasAmount = true
				result = append(result, Amount)
			} else if !hasChange {
				hasChange = true
				result = append(result, Change)
			} else {
				result = append(result, None)
			}
			continue
		}
		if !hasName {
			hasName = true
			result = append(result, Name)
		} else {
			result = append(result, None)
		}
	}
	return result
}

var datePattern = regexp.MustCompile("\\d{4}")
var numberPattern = regexp.MustCompile("-?\\d+")

func (t ImportColumnType) Validate(value string) error {
	switch t {
	case None:
		return nil
	case Name:
		if value == "" {
			return errors.New("must not be empty")
		}
		return nil
	case Date:
		if !datePattern.MatchString(value) {
			return errors.New("must match YYYY")
		}
		return nil
	case Amount:
		if !numberPattern.MatchString(value) {
			return errors.New("must be numeric")
		}
		return nil
	case Change:
		if !numberPattern.MatchString(value) {
			return errors.New("must be numeric")
		}
		return nil
	default:
		return fmt.Errorf("Unknown type %s", string(t))
	}
}

func validateRows(rows []ImportRow, name string, date string) error {
	for _, row := range rows {
		hasName := name != ""
		hasDate := date != ""
		hasAmount := false
		hasChange := false
		hasError := false
		for _, column := range row.Columns {
			switch column.Type {
			case Name:
				hasName = true
			case Date:
				hasDate = true
			case Amount:
				hasAmount = true
			case Change:
				hasChange = true
			}
			if column.Error != nil {
				hasError = true
			}
		}
		if !hasName {
			return errors.New("Name is required")
		}
		if !hasDate {
			return errors.New("Date is required")
		}
		if !hasAmount && !hasChange {
			return errors.New("Valid Amount or Change")
		}
		if hasError {
			return errors.New("Invalid entry")
		}
	}
	return nil
}
