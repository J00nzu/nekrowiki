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

func NewSessionData() *sessionData {
	sessData := sessionData{}
	sessData.lastUsed = time.Now()
	sessData.authenticated = false
	return &sessData
}

var sessions map[string]*sessionData = make(map[string]*sessionData)
var ssMutex = sync.Mutex{}
var lastSessionStoreClean = time.Now()

const sessionExpirationTime = time.Minute
const sessionStoreCleanInterval = time.Minute/2

func createSessionCookie(w http.ResponseWriter) *sessionData{

	sesID := uniuri.NewLen(32)

	sessionCookie := http.Cookie{Name: "session", Value: sesID}

	data := NewSessionData()

	ssMutex.Lock()
	sessions[sesID] = data
	ssMutex.Unlock()


	http.SetCookie(w, &sessionCookie)

	return data
}


func cleanSessionStore(){

	if lastSessionStoreClean.Add(sessionStoreCleanInterval).After(time.Now()) {
		return
	}

	lastSessionStoreClean = time.Now()

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

	ssMutex.Unlock()
}

func getSessionCookie (w http.ResponseWriter, r *http.Request) *sessionData {
	sessionCookie, err := r.Cookie("session")

	if(err != nil){

		return createSessionCookie(w)

	}else{
		if val, ok := isSessionValid(sessionCookie.Value); ok{
			return val
		}else{
			log.Println("Invalid session. Issuing new session cookie")
			return createSessionCookie(w)
		}
	}
}

func isSessionValid(session string) (*sessionData, bool) {
	cleanSessionStore()

	ssMutex.Lock()
	returnVal := &sessionData{authenticated:false}
	valid := false

	if val, ok := sessions[session]; ok {
		if val.lastUsed.Add(sessionExpirationTime).After(time.Now()) {

			val.lastUsed = time.Now() //Update last used time
			returnVal = val
			valid = true

		}else{
			log.Printf("Session %s has expired.", session)
		}
	}else{
		log.Printf("Session %s not found.", session)
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
	log.Print("login request")
	log.Print(r)

	switch r.Method {
	//GET displays the upload form.
	case "GET":
		token, err := r.Cookie("token")
		if(err == nil) {

			if(token.Value=="Authorized"){

				getSessionCookie(w,r).authenticated = true

				redirect, err := r.Cookie("afterLogin")
				//Redirect out of login
				if(err == nil){
					redirect.Expires = time.Now().Add(- time.Hour*24)
					r.AddCookie(redirect) //Expire the cookie immediately

					http.Redirect(w,r, redirect.Value, http.StatusTemporaryRedirect)
					log.Printf("redirecting to: %s ", redirect.Value)
				}else{
					http.Redirect(w,r, "/", http.StatusTemporaryRedirect)
					log.Printf("redirecting to index.html")
				}

				return;
			}
		}
		log.Print("Serving html/login.html")

		http.ServeFile(w, r, "html/login.html")

	case "POST":
		log.Print("login POST")
		err := r.ParseForm()
		if(err!=nil){
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		form := r.Form
		user := form.Get("username")
		pass := form.Get("password")

		if(user == "user" && pass == "pass"){

			getSessionCookie(w,r).authenticated = true

			redirect, err := r.Cookie("afterLogin")
			//Redirect out of login
			if(err == nil){
				redirect.Expires = time.Now().Add(- time.Hour*24)
				r.AddCookie(redirect) //Expire the cookie immediately

				http.Redirect(w,r, redirect.Value, http.StatusTemporaryRedirect)
				log.Printf("redirecting to: %s ", redirect.Value)
			}else{
				http.Redirect(w,r, "/", http.StatusTemporaryRedirect)
				log.Printf("redirecting to index.html")
			}

		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

//Authentication Middleware
func authMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie := getSessionCookie(w,r)

		if(cookie.authenticated){
			log.Println("Authenticated")
			next.ServeHTTP(w, r)
		}else{
			log.Println("Unauthenticated")
			redirectLogin(w, r)
		}

	})
}
