package main

import (
	"net/http"
	"log"
	"time"
	"github.com/dchest/uniuri"
	"sync"
)

type sessionData struct{
	lastUsed time.Time
	authenticated bool
}

func NewSessionData() sessionData {
	sessData := sessionData{}
	sessData.lastUsed = time.Now()
	sessData.authenticated = false
	return sessData
}

var sessions map[string]sessionData = make(map[string]sessionData)
var ssMutex = sync.Mutex{}
var lastSessionStoreClean = time.Now()

const sessionExpirationTime = time.Minute
const sessionStoreCleanInterval = time.Minute/2

func createSessionCookie(w http.ResponseWriter) http.Cookie{

	sesID := uniuri.NewLen(32)

	sessionCookie := http.Cookie{Name: "session", Value: sesID}

	ssMutex.Lock()
	sessions[sesID] = NewSessionData()
	ssMutex.Unlock()


	http.SetCookie(w, &sessionCookie)

	return sessionCookie
}

func cleanSessionStore(){

	if lastSessionStoreClean.Add(sessionStoreCleanInterval).After(time.Now()) {
		return
	}

	ssMutex.Lock()

	log.Println("cleaning session store...")
	log.Println(sessions)

	for k, v := range sessions {
		if v.lastUsed.Add(sessionExpirationTime).Before(time.Now()) {
			delete(sessions, k)
		}
	}

	log.Println("cleanup complete...")
	log.Println(sessions)

	lastSessionStoreClean = time.Now()
	ssMutex.Unlock()
}

func isSessionValid(session string) (sessionData, bool) {
	cleanSessionStore()

	ssMutex.Lock()
	returnVal := sessionData{authenticated:false}
	valid := false

	if val, ok := sessions[session]; ok {
		if val.lastUsed.Add(sessionExpirationTime).After(time.Now()) {

			val.lastUsed = time.Now() //Update last used time
			returnVal = val
			valid = true

		}else{
			log.Printf("Session %s has expired. creating blank session data.", session)

			sessions[session] = NewSessionData()
		}
	}else{
		log.Printf("Session %s not found. creating blank session data.", session)
		sessions[session] = NewSessionData()
	}

	ssMutex.Unlock()
	return returnVal, valid
}

func redirectLogin(w http.ResponseWriter, r *http.Request){
	//set redirect after login
	redirectAfterLogin := http.Cookie{Name: "afterLogin", Value: r.URL.Path}
	http.SetCookie(w, &redirectAfterLogin)

	log.Print("redirecting to: "+"/login ")

	http.Redirect(w,r, "/login", http.StatusTemporaryRedirect)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET displays the upload form.
	case "GET":
		token, err := r.Cookie("token")
		if(err == nil) {

			if(token.Value=="Authorized"){

				redirect, err := r.Cookie("afterLogin")
				//Redirect out of login
				if(err == nil){
					http.Redirect(w,r, redirect.Value, http.StatusTemporaryRedirect)
				}else{
					http.Redirect(w,r, "/", http.StatusTemporaryRedirect)
				}

				return;
			}
		}

		http.ServeFile(w, r, "html/login.html")
	case "POST":

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

//Authentication Middleware
func authMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")

		if(err != nil){
			log.Println(err)

			//set session cookie
			createSessionCookie(w)

			redirectLogin(w, r)

			return
		}

		log.Println(sessionCookie)

		if val, ok := isSessionValid(sessionCookie.Value); ok && val.authenticated{
			log.Println("Authenticated")
			next.ServeHTTP(w, r)
		}else{
			log.Println("Unauthenticated")

			redirectLogin(w, r)
		}

	})
}
