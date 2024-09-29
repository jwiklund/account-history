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
	"github.com/jwiklund/ah/csv"
	"github.com/jwiklund/ah/history"
	"github.com/jwiklund/ah/view"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Assets   string `default:"" usage:"Assets directory"`
	Accounts string `default:"accounts.txt" usage:"Accounts storage file"`
	Port     int    `default:"8080" usage:"Listen port"`
	Plugins  string `default:"plugins.yaml" usage:"Plugins yaml"`
}

type PluginConfig struct {
	Import []csv.ImportPlugin
}

func main() {
	userDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Warning, could not determine home config dir: %v", err)
	}
	var config Config
	loader := aconfig.LoaderFor(&config, aconfig.Config{
		EnvPrefix:  "APP",
		FlagPrefix: "app",
		Files:      []string{strings.Join([]string{userDir, "account-history", "config.yaml"}, string(os.PathSeparator))},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})
	flagSet := loader.Flags()
	initHistory := flagSet.Bool("init-history", false, "Initialize history if it does not exist")
	if err := loader.Load(); err != nil {
		panic(err)
	}
	pluginConfig, err := loadPluginConfig(userDir, config.Plugins)
	if err != nil {
		log.Fatal(err)
		return
	}
	serve(config, pluginConfig, *initHistory)
}

func loadPluginConfig(userDir, path string) (PluginConfig, error) {
	var config PluginConfig
	if path == "" {
		return config, nil
	}
	var pluginPath string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		pluginPath = strings.Join([]string{userDir, "account-history", path}, string(os.PathSeparator))
		if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
			if path == "plugins.yaml" {
				return config, nil
			}
			return config, fmt.Errorf("Plugin file does not exist: %w", err)
		}
	} else {
		pluginPath = path
	}
	pluginFile, err := os.Open(pluginPath)
	if err != nil {
		return config, err
	}
	defer pluginFile.Close()
	decoder := yaml.NewDecoder(pluginFile)
	decoder.KnownFields(true)
	err = decoder.Decode(&config)
	return config, err
}

func serve(config Config, plugins PluginConfig, initHistory bool) {
	accounts, err := history.Load(config.Accounts, initHistory)
	if err != nil {
		fmt.Printf("could not load history %v\n", err)
		return
	}
	renderer, err := view.New(config.Assets)
	if err != nil {
		fmt.Printf("could not initialize view %v\n", err)
		return
	}
	importPlugins := make(map[string]csv.ImportPlugin)
	for _, plugin := range plugins.Import {
		importPlugins[plugin.Name] = plugin
	}
	controller := control.Control{
		AccountsPath:  config.Accounts,
		Accounts:      accounts,
		Renderer:      renderer,
		ImportPlugins: importPlugins,
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
	router.POST("/import/plugin", controller.PrepareImportPlugin)
	router.POST("/import/column/:columnId", controller.PrepareImportColumn)
	router.POST("/import", controller.ImportData)

	fmt.Printf("Listening on http://localhost:%d", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), router))
}
