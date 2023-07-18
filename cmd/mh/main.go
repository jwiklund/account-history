package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cristalhq/aconfig"
	"github.com/jwiklund/money-history/history"
	"github.com/jwiklund/money-history/view"
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
	initAccounts := flagSet.Bool("initAccounts", false, "Initialize history if it does not exist")
	if err := loader.Load(); err != nil {
		panic(err)
	}
	accounts, err := history.Load(cfg.Accounts, *initAccounts)
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
	http.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				fmt.Printf("Could not parse form: %v", err)
			}
			date := accounts.CurrentDate()
			for key, value := range r.Form {
				if len(value) != 1 {
					continue
				}
				intValue, err := strconv.Atoi(value[0])
				if err != nil {
					fmt.Printf("Could not parse %s: %v\n", value[0], err)
					continue
				}
				if strings.HasSuffix(r.URL.Path, "/amount") {
					err = accounts.UpdateAmountBySlug(key, date, intValue)
					if err != nil {
						fmt.Printf("Could not update amount: %v", err)
					}
				}
				if strings.HasSuffix(r.URL.Path, "/change") {
					err = accounts.UpdateChangeBySlug(key, date, intValue)
					if err != nil {
						fmt.Printf("Could not update amount: %v", err)
					}
				}
			}
		}
		if err := renderer.Render(templateName("edit", r), w, accounts.Current()); err != nil {
			fmt.Fprintf(w, "Couild not render edit: %v", err)
		}
	})
	http.HandleFunc("/save/", func(w http.ResponseWriter, r *http.Request) {
		accounts.Save(cfg.Accounts)
		if err := renderer.Render(templateName("index", r), w, accounts.Summary()); err != nil {
			fmt.Fprintf(w, "Could not render index: %v", err)
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
