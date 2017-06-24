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
	"strings"
    "os"
	//"html/template"
	//"bytes"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET displays the upload form.
	case "GET":
		http.NotFound(w, r)
		
	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		//get the multipart reader for the request.
		reader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			//if part.FileName() is empty, skip this iteration.
			if part.FileName() == "" {
				continue
			}
			dst, err := os.Create("./uploads/" + part.FileName())
			defer dst.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			log.Println("File Uploaded: "+part.FileName());
		}
		
		 io.WriteString(w, "success");
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}


func main() {

	wikibase_b, err := ioutil.ReadFile("html/wikipage.html")
	
	
	if(err != nil){
		fmt.Print("Can't load html/wikipage.html\n")
		fmt.Print(err)
		fmt.Print("\n")
		return
	}
	
	wikibase := string(wikibase_b)
	
	fs := http.FileServer(http.Dir("public_html"))
	http.Handle("/", fs)
	
	http.HandleFunc("/upload", uploadHandler)


    http.HandleFunc("/md/", func(w http.ResponseWriter, request *http.Request) {
		url := strings.Replace(request.URL.Path, "/md", "", 1)
		
		markdown, err := ioutil.ReadFile("markdown"+ url +".md")
		
		if err != nil {
			io.WriteString(w, "an error has occurred.")
			fmt.Print(err)
			fmt.Print("\n")
		}else{
		
			complete := strings.Replace(wikibase, "%MARKDOWN%", string(github_flavored_markdown.Markdown(markdown)), -1)
			complete = strings.Replace(complete, "%NAME%", url, -1)

			w.Write([]byte(complete))
			
			/*
			io.WriteString(w, `<html><head><meta charset="utf-8"><link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
			w.Write(github_flavored_markdown.Markdown(markdown))
			io.WriteString(w, `</article></body></html>`)
			*/
		}
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
	
	log.Println("Listening...")
    log.Fatal(http.ListenAndServe(":8081", nil))

}