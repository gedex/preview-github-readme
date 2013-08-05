// Preview your local GitHub README.md.
// @author Akeda Bagus <admin@gedex.web.id>
//
// Licensed under The MIT License
// Redistributions of files must retain the above copyright notice.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

const (
	ua        = "preview-github-readme.go v0.1.0"
	url       = "https://api.github.com/markdown/raw"
	importStr = "github.com/gedex/preview-github-readme"
)

var (
	css            string
	serve          string
	serveValidator = regexp.MustCompile("^[0-9]+$")
	tpl            *template.Template
)

// getRenderedReadme reads from filename and returns
// HTML string of the Markdown's file
func getRenderedReadme(filename string) (string, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	html, err := getParsedMarkdown(buf)
	if err != nil {
		return "", err
	}

	cssPath := "./style.css"
	if pkg, err := build.Import(importStr, "", build.FindOnly); err == nil {
		p := filepath.Join(pkg.Dir, "style.css")
		if _, err := os.Stat(p); err == nil {
			cssPath = p
		}
	}
	css, err := ioutil.ReadFile(cssPath)
	if err != nil {
		return "", nil
	}

	var out = bytes.NewBuffer(nil)
	var tplData = &struct{ Css, Html string }{
		fmt.Sprintf("<style type=\"text/css\">%s</style>", css),
		html}

	err = tpl.ExecuteTemplate(out, "template.html", tplData)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// getParsedMarkdown makes a POST to GitHub endpoint to render given bytes
// and returns it as HTML string
func getParsedMarkdown(buf []byte) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("User-Agent", ua)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// usage prints the usage
func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()

	os.Exit(1)
}

func init() {
	tplPath := "./template.html"

	if pkg, err := build.Import(importStr, "", build.FindOnly); err == nil {
		p := filepath.Join(pkg.Dir, "template.html")
		if _, err := os.Stat(p); err == nil {
			tplPath = p
		}
	}

	tpl = template.Must(template.ParseFiles(tplPath))
}

func main() {

	// Set and parse flags
	flag.StringVar(&serve, "serve", "", "Serves a webserver with specified port")
	flag.Parse()

	// Renders Markdown file, passed in last arg, into HTML string
	html, err := getRenderedReadme(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)

		os.Exit(1)
	}

	if serve == "" { // Use Stdout
		fmt.Println(html)
	} else { // Listen on given `serve` argument
		if !serveValidator.MatchString(serve) {
			usage()
		}
		fmt.Println("Listening on port", serve)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, html)
		})
		http.ListenAndServe(":"+serve, nil)
	}
}
