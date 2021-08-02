package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func render(filename string, data interface{}) (string, error) {
	tpl, err := template.New("").Funcs(map[string]interface{}{
		"page": func() string {
			return filename
		},
		"year": func() string {
			return time.Now().Format("2006")
		},
	}).ParseFiles(filename)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.ExecuteTemplate(&buf, filepath.Base(filename), data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderMain(filename string, data interface{}) (string, error) {
	body, err := render(filename, data)
	if err != nil {
		return "", err
	}
	return render("books/layout.gohtml", struct{ Body template.HTML }{
		Body: template.HTML(body),
	})
}

func renderIntro(filename string, data interface{}) (string, error) {
	body, err := render(filename, data)
	if err != nil {
		return "", err
	}
	return renderMain("books/intro/layout.gohtml", struct {
		Body template.HTML
		Page string
	}{
		Body: template.HTML(body),
		Page: filename,
	})
}

var router *httprouter.Router

type Redirect struct {
	From, To string
	Code     int
}

type Generator struct {
	redirects []Redirect
}

func (g *Generator) generatePage(path string, tpl Template) error {
	str, err := tpl.Render(nil)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	p := "./build" + path
	if strings.HasSuffix(p, "/") {
		p += "index.html"
	} else {
		g.redirects = append(g.redirects, Redirect{From: p + "/", To: p, Code: http.StatusPermanentRedirect})
		p += ".html"
	}

	err = os.MkdirAll(filepath.Dir(p), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.WriteString(f, str)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateRedirects() error {
	f, err := os.Create("./build/_redirects")
	if err != nil {
		return err
	}
	defer f.Close()

	for _, r := range g.redirects {
		_, err = fmt.Fprintf(f, "%s %s %d\n", r.From, r.To, r.Code)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateHTML(path string, tpl Template) {
	f := func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		str, err := tpl.Render(nil)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(res, str)
	}
	router.GET(path, f)
	if strings.HasSuffix(path, "/") {
		router.GET(path+"index.htm", f)
		router.GET(path+"index.html", f)
	} else {
		router.GET(path+".htm", f)
		router.GET(path+".html", f)
		router.GET(path+"/", f)
		router.GET(path+"/index.htm", f)
		router.GET(path+"/index.html", f)
	}
	{
		str, err := tpl.Render(nil)
		if errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			panic(err)
		}

		p := "./build" + path
		if strings.HasSuffix(p, "/") {
			p += "index.html"
		} else {
			p += ".html"
		}

		err = os.MkdirAll(filepath.Dir(p), 0755)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(p)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = io.WriteString(f, str)
		if err != nil {
			panic(err)
		}
	}
}

func build() error {
	_ = os.RemoveAll("./build")

	g := new(Generator)

	type Page struct {
		From     string
		Template Template
	}
	var pages []Page
	pages = append(pages, []Page{
		{"/", PageTemplate{FileTemplate("index.gohtml"), ""}},
		{"/guides/machine_setup", guideTemplate{FileTemplate("guides/01_machine_setup.gohtml"), "Machine Setup"}},
		{"/guides/bootcamp", guideTemplate{FileTemplate("guides/02_bootcamp.gohtml"), "Bootcamp"}},
	}...)
	for w := 1; w <= 4; w++ {
		for d := 1; d <= 5; d++ {
			pages = append(pages, Page{
				fmt.Sprintf("/guides/bootcamp/week-%d/day-%d", w, d),
				guideTemplate{
					FileTemplate(fmt.Sprintf("guides/bootcamp/week-%d/day-%d.gohtml", w, d)),
					"Bootcamp",
				},
			})
		}
	}

	pages = append(pages, Page{"/books/intro", introTemplate{Template: FileTemplate("books/intro/front.gohtml")}})
	sections := []string{
		"Getting Started",
		"Your First Program",
		"Types",
		"Variables",
		"Control Structures",
		"Arrays, Slices and Maps",
		"Functions",
		"Pointers",
		"Structs and Interfaces",
		"Concurrency",
		"Packages",
		"Testing",
		"The Core Packages",
		"Next Steps",
	}
	for i, section := range sections {
		pages = append(pages, Page{
			fmt.Sprint("/books/intro/", i+1),
			introTemplate{
				Template: FileTemplate(fmt.Sprint("books/intro/", i+1, ".gohtml")),
				Title:    section,
			},
		})
	}
	for i := 1; i <= 14; i++ {
		for _, str := range []string{
			fmt.Sprintf("/%d", i),
			fmt.Sprintf("/%d/index.htm", i),
		} {
			g.redirects = append(g.redirects, Redirect{
				From: str,
				To:   fmt.Sprintf("/books/intro/%d", i),
				Code: http.StatusPermanentRedirect,
			})
		}
	}

	pages = append(pages, Page{"/books/web", webTemplate{Template: FileTemplate("books/web/front.gohtml")}})
	pages = append(pages, Page{"/books/web/00-01", webTemplate{Template: FileTemplate("books/web/00-01.gohtml")}})
	pages = append(pages, Page{"/books/web/01-01", webTemplate{Template: FileTemplate("books/web/01-01.gohtml")}})
	pages = append(pages, Page{"/books/web/01-02", webTemplate{Template: FileTemplate("books/web/01-02.gohtml")}})
	pages = append(pages, Page{"/books/web/01-03", webTemplate{Template: FileTemplate("books/web/01-03.gohtml")}})
	pages = append(pages, Page{"/books/web/01-04", webTemplate{Template: FileTemplate("books/web/01-04.gohtml")}})

	for i := 1; i < 10; i++ {
		pages = append(pages, Page{
			fmt.Sprintf("/books/web/02-%02d", i),
			webTemplate{Template: FileTemplate(fmt.Sprintf("books/web/02-%02d.gohtml", i))},
		})
	}

	for _, page := range pages {
		err := g.generatePage(page.From, page.Template)
		if err != nil {
			return err
		}
	}

	err := g.generateRedirects()
	if err != nil {
		return err
	}

	err = filepath.Walk("./public", func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		dstPath := filepath.Join("./build", p)
		err = os.MkdirAll(filepath.Dir(dstPath), 0755)
		if err != nil {
			return err
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		src, err := os.Open(p)
		if err != nil {
			return err
		}
		defer src.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
