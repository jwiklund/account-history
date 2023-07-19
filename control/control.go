package control

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/jwiklund/ah/history"
	"github.com/jwiklund/ah/view"
)

type Control struct {
	AccountsPath string
	Accounts     *history.Accounts
	Renderer     view.Renderer
}

func templateName(part string, r *http.Request) string {
	hx, _ := r.Header["Hx-Request"]
	suffix := "html"
	if len(hx) == 1 && hx[0] == "true" {
		suffix = "body.html"
	}
	return strings.Join([]string{part, suffix}, ".")
}

func slugAndIntValue(r *http.Request) (string, int, error) {
	for key, value := range r.Form {
		if len(value) != 1 {
			continue
		}
		intValue, err := strconv.Atoi(value[0])
		if err != nil {
			return "", 0, err
		}
		return key, intValue, nil
	}
	return "", 0, errors.New("no form values")
}
