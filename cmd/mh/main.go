package main

import (
	"fmt"
	"net/http"

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
		if err := renderer.Render("index.html", w, accounts.Summary()); err != nil {
			fmt.Fprintf(w, "Could not render index: %v", err)
		}
	})
	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
