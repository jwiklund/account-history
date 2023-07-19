package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cristalhq/aconfig"
	"github.com/jwiklund/ah/history"
	"github.com/jwiklund/ah/view"
)

type Config struct {
	Assets   string `default:"" usage:"Assets directory"`
	Accounts string `default:"accounts.txt" usage:"Accounts storage file"`
	Port     int    `default:"8080" usage:"Listen port"`
}

func main() {
	var cfg Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		EnvPrefix:  "APP",
		FlagPrefix: "app",
	})
	flagSet := loader.Flags()
	initHistory := flagSet.Bool("init-history", false, "Initialize history if it does not exist")
	if err := loader.Load(); err != nil {
		panic(err)
	}
	accounts, err := history.Load(cfg.Accounts, *initHistory)
	if err != nil {
		fmt.Printf("could not load history %v\n", err)
		return
	}
	renderer, err := view.New(cfg.Assets)
	if err != nil {
		fmt.Printf("could not initialize view %v\n", err)
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := renderer.Render(templateName("index", r), w, accounts.Summary()); err != nil {
			fmt.Fprintf(w, "Could not render index: %v", err)
		}
	})
	http.HandleFunc("/save/", func(w http.ResponseWriter, r *http.Request) {
		accounts.Save(cfg.Accounts)
		if err := renderer.Render(templateName("index", r), w, accounts.Summary()); err != nil {
			fmt.Fprintf(w, "Could not render index: %v", err)
		}
	})
	http.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				fmt.Printf("Could not parse form: %v", err)
			}
			date := accounts.CurrentDate()
			if strings.HasSuffix(r.URL.Path, "/amount") {
				key, value, err := slugAndIntValue(r)
				if err != nil {
					fmt.Printf("Could not parse amount: %v", err)
				} else {
					err = accounts.UpdateAmountBySlug(key, date, value)
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
					err = accounts.UpdateChangeBySlug(key, date, value)
					if err != nil {
						fmt.Printf("Could not update amount: %v", err)
					}
				}
			}
			if strings.HasSuffix(r.URL.Path, "/add") {
				name := r.Form.Get("name")
				if name != "" {
					accounts.AddAccount(name, date)
				}
			}
		}
		if err := renderer.Render(templateName("edit", r), w, accounts.Current()); err != nil {
			fmt.Fprintf(w, "Couild not render edit: %v", err)
		}
	})
	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
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
