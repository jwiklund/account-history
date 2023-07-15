package main

import (
	"fmt"
	"net/http"

	"github.com/cristalhq/aconfig"
	"github.com/jwiklund/money-history/view"
)

type Config struct {
	Assets string `default:"" usage:"Assets directory"`
	Port   int    `default:"8080" usage:"Listen port"`
}

func main() {
	var cfg Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		EnvPrefix:  "APP",
		FlagPrefix: "app",
	})
	if err := loader.Load(); err != nil {
		panic(err)
	}
	renderer, err := view.New(cfg.Assets)
	if err != nil {
		fmt.Printf("could not initialize view %v", err)
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := renderer.Render("index.html", w, nil); err != nil {
			fmt.Fprintf(w, "Could not render index: %v", err)
		}
	})
	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
