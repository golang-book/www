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

	"github.com/badgerodon/socketmaster/client"
	"github.com/badgerodon/socketmaster/protocol"
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
	assets := http.FileServer(http.Dir("assets"))

	router = httprouter.New()
	register("/", PageTemplate{FileTemplate("index.gohtml")})
	register("/guides/machine_setup", PageTemplate{FileTemplate("guides/01_machine_setup.gohtml")})
	register("/books/intro", IntroTemplate{Template: FileTemplate("books/intro/front.gohtml")})
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
			IntroTemplate{
				Template: FileTemplate(fmt.Sprint("books/intro/", i+1, ".gohtml")),
				Title:    section,
			},
		)
	}
	router.GET("/assets/*path", func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		req.URL.Path = "/" + params.ByName("path")
		assets.ServeHTTP(res, req)
	})
	handler := router

	log.Println("starting server on :443")
	li1, err := client.Listen(protocol.SocketDefinition{
		Port: 443,
		HTTP: &protocol.SocketHTTPDefinition{
			DomainSuffix: "golang-book.com",
		},
		TLS: &protocol.SocketTLSDefinition{
			Cert: tlsCert,
			Key:  tlsKey,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer li1.Close()
	go http.Serve(li1, handler)

	log.Println("starting server on :80")
	li2, err := client.Listen(protocol.SocketDefinition{
		Port: 80,
		HTTP: &protocol.SocketHTTPDefinition{
			DomainSuffix: "golang-book.com",
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer li2.Close()

	err = http.Serve(li2, handler)
	if err != nil {
		log.Fatalln(err)
	}
}
