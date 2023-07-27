package control

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jwiklund/ah/csv"
	"github.com/jwiklund/ah/history"
	"github.com/jwiklund/ah/view"
)

type Control struct {
	AccountsPath  string
	Accounts      *history.Accounts
	Renderer      view.Renderer
	ImportPlugins map[string]csv.ImportPlugin
}

func templateName(base string, r *http.Request) string {
	return templateNameWithPart(base, r, "body")
}

func templateNameWithPart(base string, r *http.Request, part string) string {
	parts := []string{base}
	if isHx(r) {
		parts = append(parts, part)
	}
	parts = append(parts, "html")
	return strings.Join(parts, ".")
}

func isHx(r *http.Request) bool {
	hx, _ := r.Header["Hx-Request"]
	return len(hx) == 1 && hx[0] == "true"
}

func formInput(r *http.Request, key string) string {
	return strings.TrimSpace(rawFormInput(r, key))
}

func formIntInput(r *http.Request, key string) (int, error) {
	input := formInput(r, key)
	if input == "" {
		return 0, fmt.Errorf("no value for %s", key)
	}
	input = strings.ReplaceAll(input, ",", "")
	input = strings.ReplaceAll(input, " ", "")
	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %v", key, err)
	}
	return value, nil
}

func rawFormInput(r *http.Request, key string) string {
	if len(r.Form[key]) == 1 {
		return r.Form[key][0]
	}
	return ""
}

func slugAndIntValue(r *http.Request) (string, int, error) {
	for key, value := range r.Form {
		if len(value) != 1 {
			continue
		}
		input := value[0]
		input = strings.ReplaceAll(input, ",", "")
		input = strings.ReplaceAll(input, " ", "")
		intValue, err := strconv.Atoi(input)
		if err != nil {
			return "", 0, err
		}
		return key, intValue, nil
	}
	return "", 0, errors.New("no form values")
}
