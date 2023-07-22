package csv

import (
	"testing"

	"github.com/jwiklund/ah/csv/assets"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestImport(t *testing.T) {
	testFiles, err := assets.EmbedFS.ReadDir(".")
	assert.NoError(t, err)
	for _, testFilename := range testFiles {
		testFile, err := assets.EmbedFS.Open(testFilename.Name())
		assert.NoError(t, err)
		var testCase ImportTestCase
		err = yaml.NewDecoder(testFile).Decode(&testCase)
		assert.NoError(t, err, testFilename.Name())
		if err != nil {
			continue
		}

		name := testFilename.Name()[0 : len(testFilename.Name())-5]
		t.Run(name, func(t2 *testing.T) {
			rows, columns, name, date, err := Import(testCase.Input, testCase.Opts)
			if testCase.Error == "" {
				assert.NoError(t2, err)
			} else {
				assert.Error(t2, err, testCase.Error)
			}
			assert.Equal(t2, testCase.Expect.Rows, rows.Rows)
			assert.Equal(t2, testCase.Expect.Columns, columns)
			assert.Equal(t2, testCase.Expect.Name, name)
			assert.Equal(t2, testCase.Expect.Date, date)
		})
	}
}

type ImportTestCase struct {
	Input  []string
	Opts   ImportOptions
	Expect struct {
		Rows    []ImportRow
		Columns []ImportColumnType
		Name    string
		Date    string
	}
	Error string
}
