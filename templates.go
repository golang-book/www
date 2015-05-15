package main

import (
	"bytes"
	"html/template"
	"path/filepath"
	"time"
)

type Template interface {
	Render(data interface{}) (string, error)
}

type FileTemplate string

func (ft FileTemplate) Render(data interface{}) (string, error) {
	tpl, err := template.New("").Funcs(map[string]interface{}{
		"year": func() string {
			return time.Now().Format("2006")
		},
	}).ParseFiles("templates/" + string(ft))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.ExecuteTemplate(&buf, filepath.Base(string(ft)), data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type IntroTemplate struct {
	Template
	Title string
}

func (it IntroTemplate) Render(data interface{}) (string, error) {
	body, err := it.Template.Render(data)
	if err != nil {
		return "", err
	}
	return PageTemplate{FileTemplate("books/intro/layout.gohtml")}.Render(struct {
		Body  template.HTML
		Title string
	}{
		Body:  template.HTML(body),
		Title: it.Title,
	})
}

type PageTemplate struct {
	Template
}

func (pt PageTemplate) Render(data interface{}) (string, error) {
	body, err := pt.Template.Render(data)
	if err != nil {
		return "", err
	}
	return FileTemplate("layout.gohtml").Render(struct {
		Body template.HTML
	}{
		Body: template.HTML(body),
	})
}
