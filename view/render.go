package view

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"

	"github.com/jwiklund/ah/view/assets"
)

type Renderer interface {
	Render(name string, wr io.Writer, data any) error
}

type CheckRenderer interface {
	Renderer
	Check() error
}

func New(assetsPath string) (Renderer, error) {
	if assetsPath == "" {
		templates, err := template.New("").Funcs(map[string]any{
			"json": toJson,
		}).ParseFS(assets.EmbedFs, "*.html")
		if err != nil {
			return nil, fmt.Errorf("could not parse templates: %w", err)
		}
		renderer := eagerRenderer{templates}
		return renderer, check(renderer)
	}
	renderer := lazyRenderer{os.DirFS(assetsPath)}
	return renderer, check(renderer)
}

func check(r Renderer) error {
	return r.Render("index.html", io.Discard, nil)
}

type eagerRenderer struct {
	templates *template.Template
}

func (e eagerRenderer) Render(name string, wr io.Writer, data any) error {
	return e.templates.ExecuteTemplate(wr, name, data)
}

type lazyRenderer struct {
	fs fs.FS
}

func (l lazyRenderer) Render(name string, wr io.Writer, data any) error {
	templates, err := template.New("").Funcs(map[string]any{
		"json": toJson,
	}).ParseFS(l.fs, "*.html")
	if err != nil {
		return fmt.Errorf("could not parse templates: %w", err)
	}
	err = templates.ExecuteTemplate(wr, name, data)
	return err
}

func toJson(data any) (template.JS, error) {
	result, err := json.MarshalIndent(data, "", "  ")
	return template.JS(result), err
}
