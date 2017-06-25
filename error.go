package main

import (
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"log"
	"strings"
	"strconv"
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
