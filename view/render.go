package view

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/jwiklund/ah/view/assets"
)

type Renderer interface {
	Render(name string, wr io.Writer, data any) error
	Resource(name string, wr io.Writer) error
}

type CheckRenderer interface {
	Renderer
	Check() error
}

var funcMap = map[string]any{
	"json":  toJson,
	"human": toHuman,
}

func New(assetsPath string) (Renderer, error) {
	if assetsPath == "" {
		templates, err := template.New("").Funcs(funcMap).ParseFS(assets.EmbedFs, "*.html")
		if err != nil {
			return nil, fmt.Errorf("could not parse templates: %w", err)
		}
		renderer := eagerRenderer{templates, assets.EmbedFs}
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
	fs        fs.FS
}

func (e eagerRenderer) Render(name string, wr io.Writer, data any) error {
	return e.templates.ExecuteTemplate(wr, name, data)
}

func (e eagerRenderer) Resource(name string, wr io.Writer) error {
	return resource(e.fs, name, wr)
}

type lazyRenderer struct {
	fs fs.FS
}

func (l lazyRenderer) Resource(name string, wr io.Writer) error {
	return resource(l.fs, name, wr)
}

func resource(fs fs.FS, name string, wr io.Writer) error {
	f, err := fs.Open(name)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", name, err)
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	_, err = wr.Write(data)
	return err
}

func (l lazyRenderer) Render(name string, wr io.Writer, data any) error {
	templates, err := template.New("").Funcs(funcMap).ParseFS(l.fs, "*.html")
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

func toHuman(data any) string {
	if i, ok := data.(int); ok {
		return humanize.Comma(int64(i))
	}
	return fmt.Sprintf("%v", data)
}
