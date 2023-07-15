package view

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"text/template"
)

type Renderer interface {
	Render(name string, wr io.Writer, data any) error
}

func New(assets string) (Renderer, error) {
	if assets == "" {
		return nil, errors.New("Embedded assets not implemented yet")
	}
	result := DebugAssets{assets}
	err := result.check()
	return result, err
}

type DebugAssets struct {
	assets string
}

func (d DebugAssets) check() error {
	return d.Render("index.html", io.Discard, nil)
}

func (d DebugAssets) Render(name string, wr io.Writer, data any) error {
	namedTemplate, err := template.ParseFiles(strings.Join([]string{d.assets, name}, "/"))
	if err != nil {
		return fmt.Errorf("could not parse %s: %w", name, err)
	}
	err = namedTemplate.Execute(wr, data)
	return err
}
