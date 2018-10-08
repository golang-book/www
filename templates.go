package main

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"
	"time"
)

type Template interface {
	Render(data interface{}) (string, error)
}

type FileTemplate string

func (ft FileTemplate) Render(data interface{}) (string, error) {
	tpl, err := template.New("").Funcs(map[string]interface{}{
		"public_link": func(name string) string {
			version := getVersion("public/" + name)
			parts := strings.SplitN(name, ".", 2)
			if len(parts) == 1 {
				return "/public/" + name + "." + version
			}
			return "/public/" + parts[0] + "." + version + "." + parts[1]
		},
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

type introTemplate struct {
	Template
	Title string
}

func (it introTemplate) Render(data interface{}) (string, error) {
	body, err := it.Template.Render(data)
	if err != nil {
		return "", err
	}
	title := "An Introduction to Programming in Go"
	if it.Title != "" {
		title = it.Title + " — " + title
	}
	return PageTemplate{
		Template: FileTemplate("books/intro/layout.gohtml"),
		Title:    title,
	}.Render(struct {
		Body  template.HTML
		Title string
	}{
		Body:  template.HTML(body),
		Title: it.Title,
	})
}

type webTemplate struct {
	Template
	Title string
}

func (wt webTemplate) Render(data interface{}) (string, error) {
	body, err := wt.Template.Render(data)
	if err != nil {
		return "", err
	}
	title := "An Introduction to Web Development in Go"
	if wt.Title != "" {
		title = wt.Title + " — " + title
	}
	return PageTemplate{
		Template: FileTemplate("books/web/layout.gohtml"),
		Title:    title,
	}.Render(struct {
		Body  template.HTML
		Title string
	}{
		Body:  template.HTML(body),
		Title: wt.Title,
	})
}

type distributedSystemsTemplate struct {
	Template
	Subtitle string
}

func (dst distributedSystemsTemplate) Render(data interface{}) (string, error) {
	paging, err := FileTemplate("books/distributed-systems/paging.gohtml").Render(struct {
		Previous, Next string
	}{
		Previous: "???",
		Next:     "???",
	})
	if err != nil {
		return "", err
	}

	body, err := dst.Template.Render(struct {
		Paging template.HTML
	}{
		Paging: template.HTML(paging),
	})
	if err != nil {
		return "", err
	}

	title := "An Introduction to Distributed Systems in Go"
	return PageTemplate{
		Template: FileTemplate("books/distributed-systems/layout.gohtml"),
		Title:    title,
	}.Render(struct {
		Body     template.HTML
		Title    string
		Subtitle string
	}{
		Body:     template.HTML(body),
		Title:    title,
		Subtitle: dst.Subtitle,
	})
}

type guideTemplate struct {
	Template
	Title string
}

func (gt guideTemplate) Render(data interface{}) (string, error) {
	body, err := gt.Template.Render(data)
	if err != nil {
		return "", err
	}
	return PageTemplate{
		Template: FileTemplate("guides/layout.gohtml"),
		Title:    gt.Title,
	}.Render(struct {
		Body template.HTML
	}{
		Body: template.HTML(body),
	})
}

type PageTemplate struct {
	Template
	Title string
}

func (pt PageTemplate) Render(data interface{}) (string, error) {
	body, err := pt.Template.Render(data)
	if err != nil {
		return "", err
	}
	return FileTemplate("layout.gohtml").Render(struct {
		Body  template.HTML
		Title string
	}{
		Body:  template.HTML(body),
		Title: pt.Title,
	})
}
