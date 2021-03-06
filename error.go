package main

import (
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"log"
	"strings"
	"strconv"
	"time"
)

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

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		//TODO: remove time logging
		t1 := time.Now()

		next.ServeHTTP(w, r)

		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}
