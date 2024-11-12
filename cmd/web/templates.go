package main

import (
	"path/filepath"
	"text/template"
)

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	partials, err := filepath.Glob("./ui/html/partials/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		files := []string{
			"./ui/html/base.html",
			page,
		}
		files = append(files, partials...)

		t, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		name := filepath.Base(page)
		cache[name] = t
	}

	return cache, nil
}
