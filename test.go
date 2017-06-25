package main

import (
	"fmt"
	"log"
	"net/http"
	/*
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
	*/
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"strings"
    "os"
	"regexp"
	"net/http/httptest"
	"strconv"
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

type ModifierMiddleware struct {
	handler http.Handler
}

func (m *ModifierMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()
	// passing a ResponseRecorder instead of the original RW
	m.handler.ServeHTTP(rec, r)
	// after this finishes, we have the response recorded
	// and can modify it before copying it to the original RW

}

func customErrorMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rec := httptest.NewRecorder()
		// passing a ResponseRecorder instead of the original RW
		next.ServeHTTP(rec, r)

		if rec.Code >= 400 { // check if code is an error
			dir, err := ioutil.ReadDir("errorpages")

			if (err == nil) {
				for _, file := range dir {
					log.Print(file.Name())
					if strings.HasPrefix(file.Name(), strconv.Itoa(rec.Code)) {

						errorRec := httptest.NewRecorder()
						http.ServeFile(errorRec, r, "errorpages/"+file.Name())

						for k, v := range errorRec.Header() {
							w.Header()[k] = v
						}
						w.WriteHeader(rec.Code) // change the code of custom page to our original error code
						w.Write(errorRec.Body.Bytes())

						return;
					}
				}
			} else {
				log.Print("Can't find folder called 'errorpages'")
			}
		}

		// we copy the original headers first
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}

		// write the correct http code
		w.WriteHeader(rec.Code)

		// then write out the original body
		w.Write(rec.Body.Bytes())

	})
}

func markdownHandler(w http.ResponseWriter, request *http.Request) {

	wikibase_b, err := ioutil.ReadFile("html/wikipage.html")

	if(err != nil){
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

func _start(args []string) {
	fs := http.FileServer(http.Dir("public_html"))
	http.Handle("/", customErrorMW(authMW(fs)))

	//http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(gfmstyle.Assets)))

	http.Handle("/login", customErrorMW(http.HandlerFunc(loginHandler)))

	ufs := http.FileServer(http.Dir("uploads"))
	http.Handle("/uploads/", customErrorMW(authMW(http.StripPrefix("/uploads", ufs))))

	http.Handle("/upload", customErrorMW(authMW(http.HandlerFunc(uploadHandler))))
	http.Handle("/md/", customErrorMW(authMW(http.StripPrefix("/md", http.HandlerFunc(markdownHandler)))))

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func _stop(args []string) {
	fmt.Print("Not implemented")
}

func _restart(args []string) {
	fmt.Print("Not implemented")
}

func _config(args []string) {
	fmt.Print("Not implemented")
}

func _help() {
	fmt.Print("Usage: \n$./nekrowiki start")
}

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		_help()
		return
	} else {
		function := args[0]

		var additional_args []string

		if len(args) > 1 {
			additional_args = args[1:]
		} else {
			additional_args = make([]string, 0)
		}

		switch function {
		case "start":
			_start(additional_args)
		case "stop":
			_stop(additional_args)
		case "restart":
			_restart(additional_args)
		case "config":
			_config(additional_args)
		default:
			_help()
		}

	}

}