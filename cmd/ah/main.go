package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/jwiklund/ah/control"
	"github.com/jwiklund/ah/history"
	"github.com/jwiklund/ah/view"
)

type Config struct {
	Assets   string `default:"" usage:"Assets directory"`
	Accounts string `default:"accounts.txt" usage:"Accounts storage file"`
	Port     int    `default:"8080" usage:"Listen port"`
}

func main() {
	userDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Warning, could not determine home config dir: %v", err)
	}
	var cfg Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		EnvPrefix:  "APP",
		FlagPrefix: "app",
		Files:      []string{strings.Join([]string{userDir, "ah", "config.yaml"}, string(os.PathSeparator))},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})
	flagSet := loader.Flags()
	initHistory := flagSet.Bool("init-history", false, "Initialize history if it does not exist")
	if err := loader.Load(); err != nil {
		panic(err)
	}
	serve(cfg, *initHistory)
}

func serve(cfg Config, initHistory bool) {
	accounts, err := history.Load(cfg.Accounts, initHistory)
	if err != nil {
		fmt.Printf("could not load history %v\n", err)
		return
	}
	renderer, err := view.New(cfg.Assets)
	if err != nil {
		fmt.Printf("could not initialize view %v\n", err)
		return
	}
	controller := control.Control{
		AccountsPath: cfg.Accounts,
		Accounts:     accounts,
		Renderer:     renderer,
	}
	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/save/", controller.Save)
	http.HandleFunc("/edit/", controller.Edit)
	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
