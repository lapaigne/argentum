package main

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type Templates struct {
	templates map[string]*template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any, c echo.Context) error {
	v, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("not found: %s", name)
	}

	if c.Request().Header.Get("HX-Target") == "content" {
		return v.ExecuteTemplate(w, "content", data)
	}
	if c.Request().Header.Get("HX-Request") == "true" {
		return v.ExecuteTemplate(w, name, data)
	}

	return v.ExecuteTemplate(w, "base", data)
}

func ParseTemplates(p string) (map[string]*template.Template, error) {

	templates := make(map[string]*template.Template)

	partials, err := template.ParseGlob(filepath.Join(p, "partials", "*.html"))
	if err != nil {
		return nil, err
	}

	pages, err := filepath.Glob(filepath.Join(p, "pages", "*.html"))
	if err != nil {
		return nil, err
	}

	for _, file := range pages {
		t, err := partials.Clone()
		if err != nil {
			return nil, err
		}

		t, err = t.ParseFiles(file)
		if err != nil {
			return nil, err
		}

		name := filepath.Base(file)
		name = name[:len(name)-len(filepath.Ext(name))]

		templates[name] = t
	}

	for _, v := range partials.Templates() {
		name := v.Name()
		templates[name] = partials
	}

	return templates, nil
}

func NewTemplates(p string) (*Templates, error) {
	templates, err := ParseTemplates(p)
	if err != nil {
		return nil, err
	}
	return &Templates{templates: templates}, nil
}
