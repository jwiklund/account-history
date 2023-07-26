package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/julienschmidt/httprouter"
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
	router := httprouter.New()
	router.GET("/", controller.Index)
	router.POST("/save", controller.Save)

	router.GET("/edit", controller.Edit)
	router.POST("/edit/add", controller.EditAdd)
	router.POST("/edit/amount", controller.EditAmount)
	router.POST("/edit/change", controller.EditChange)

	router.GET("/edit/account/:accountSlug", controller.EditAccount)
	router.POST("/edit/account/:accountSlug/add", controller.EditAccountAdd)
	router.POST("/edit/account/:accountSlug/amount/:year", controller.EditAccountAmount)
	router.POST("/edit/account/:accountSlug/change/:year", controller.EditAccountChange)

	router.GET("/import", controller.Import)
	router.POST("/import/prepare", controller.PrepareImport)
	router.POST("/import/separator", controller.PrepareImportSeparator)
	router.POST("/import/column/:columnId", controller.PrepareImportColumn)
	router.POST("/import", controller.ImportData)

	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
