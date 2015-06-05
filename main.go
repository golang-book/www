package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
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

func register(path string, tpl Template) {
	f := func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		str, err := tpl.Render(nil)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		res.Header().Set("Content-Type", "text-html; charset=utf-8")
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
}

func main() {
	log.SetFlags(0)

	router = httprouter.New()
	register("/", PageTemplate{FileTemplate("index.gohtml"), ""})
	register("/guides/machine_setup", guideTemplate{FileTemplate("guides/01_machine_setup.gohtml"), "Machine Setup"})
	register("/books/intro", introTemplate{Template: FileTemplate("books/intro/front.gohtml")})
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
		register(
			fmt.Sprint("/books/intro/", i+1),
			introTemplate{
				Template: FileTemplate(fmt.Sprint("books/intro/", i+1, ".gohtml")),
				Title:    section,
			},
		)
	}

	for i := 1; i <= 14; i++ {
		for _, str := range []string{
			fmt.Sprintf("/%d", i),
			fmt.Sprintf("/%d/index.htm", i),
		} {
			dst := fmt.Sprintf("/books/intro/%d", i)
			router.GET(str, func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
				http.Redirect(res, req, dst, http.StatusMovedPermanently)
			})
		}
	}

	register("/books/web", webTemplate{Template: FileTemplate("books/web/front.gohtml")})
	register("/books/web/00-01", webTemplate{Template: FileTemplate("books/web/00-01.gohtml")})
	register("/books/web/01-01", webTemplate{Template: FileTemplate("books/web/01-01.gohtml")})

	public := http.FileServer(http.Dir("public"))
	router.GET("/public/*path", func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		path := params.ByName("path")
		maxAge := "3600"
		parts := strings.SplitN(path, ".", 3)
		if len(parts) == 3 {
			p := parts[0] + "." + parts[2]
			if parts[1] == getVersion("public/"+p) {
				path = p
				maxAge = "31556926"
			}
		}
		req.URL.Path = "/" + path
		res.Header().Set("Cache-Control", "max-age="+maxAge+", public")
		public.ServeHTTP(res, req)
	})

	router.GET("/health", func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		io.WriteString(res, "OK")
	})

	handler := gziphandler.GzipHandler(router)

	log.Println("starting server on 127.0.0.1:8002")
	err := http.ListenAndServe("127.0.0.1:8002", handler)
	if err != nil {
		log.Fatalln(err)
	}
}
