package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
	"io"
	"io/ioutil"
)

func main() {

    http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		markdown, err := ioutil.ReadFile("markdown"+request.URL.Path+".md")
		
		if err != nil {
			io.WriteString(w, "an error has occurred.");
			fmt.Print(err)
			fmt.Print("\n");
		}else{
			io.WriteString(w, `<html><head><meta charset="utf-8"><link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
			w.Write(github_flavored_markdown.Markdown(markdown))
			io.WriteString(w, `</article></body></html>`)
		}
    })
    
    http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request){
        fmt.Fprintf(writer, "Hi %q", html.EscapeString(request.URL.Path))
    })
	
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(gfmstyle.Assets)))
	
	//markdown := []byte("# GitHub Flavored Markdown\n\nHello.")
	
	/*
	http.HandleFunc("/mdtest", func(w http.ResponseWriter, r *http.Request){
        io.WriteString(w, `<html><head><meta charset="utf-8"><link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
		w.Write(github_flavored_markdown.Markdown(markdown))
		io.WriteString(w, `</article></body></html>`)
    })
	*/
	
    log.Fatal(http.ListenAndServe(":8081", nil))

}