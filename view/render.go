package view

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"github.com/jwiklund/money-history/view/assets"
)

type Renderer interface {
	Render(name string, wr io.Writer, data any) error
}

type CheckRenderer interface {
	Renderer
	Check() error
}

func New(assets string) (Renderer, error) {
	var result CheckRenderer
	if assets == "" {
		var err error
		result, err = newEmbed()
		if err != nil {
			return nil, err
		}
	} else {
		result = DebugAssets{
			assetDir: assets,
		}
	}
	err := result.Check()
	return result, err
}

type DebugAssets struct {
	assetDir string
}

func (d DebugAssets) Check() error {
	return d.Render("index.html", io.Discard, nil)
}

func (d DebugAssets) Render(name string, wr io.Writer, data any) error {
	matches, err := filepath.Glob(strings.Join([]string{d.assetDir, "*.html"}, "/"))
	if err != nil {
		return fmt.Errorf("could not glob templates: %w", err)
	}
	templates, err := template.ParseFiles(matches...)
	if err != nil {
		return fmt.Errorf("could not parse templates: %w", err)
	}
	err = templates.ExecuteTemplate(wr, name, data)
	return err
}

type EmbedAssets struct {
	templates *template.Template
}

func newEmbed() (CheckRenderer, error) {
	templates, err := template.ParseFS(assets.EmbedFs, "*.html")
	if err != nil {
		return nil, fmt.Errorf("could not parse templates: %w", err)
	}
	return EmbedAssets{
		templates: templates,
	}, nil
}

func (e EmbedAssets) Check() error {
	return e.Render("index.html", io.Discard, nil)
}

func (e EmbedAssets) Render(name string, wr io.Writer, data any) error {
	return e.templates.ExecuteTemplate(wr, name, data)
}
