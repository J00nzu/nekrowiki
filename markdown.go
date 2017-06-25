package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"
	"regexp"
	"strings"
)

func markdownHandler(w http.ResponseWriter, request *http.Request) {

	wikibase_b, err := ioutil.ReadFile("html/wikipage.html")

	if (err != nil) {
		fmt.Print("Can't load html/wikipage.html\n")
		fmt.Print(err)
		fmt.Print("\n")
		return
	}

	wikibase := string(wikibase_b)

	//url := strings.Replace(request.URL.Path, "/md", "", 1)

	url := request.URL.Path

	markdown, err := ioutil.ReadFile("markdown" + url + ".md")

	if err != nil {
		http.NotFound(w, request)
		fmt.Print(err)
		fmt.Print("\n")
	} else {

		unsafe := blackfriday.MarkdownCommon(markdown)
		p := bluemonday.UGCPolicy()
		p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
		html := p.SanitizeBytes(unsafe)

		complete := strings.Replace(wikibase, "%MARKDOWN%", string(html), -1)
		complete = strings.Replace(complete, "%NAME%", url, -1)

		w.Write([]byte(complete))

		/*
		io.WriteString(w, `<html><head><meta charset="utf-8"><link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
		w.Write(github_flavored_markdown.Markdown(markdown))
		io.WriteString(w, `</article></body></html>`)
		*/
	}
}
