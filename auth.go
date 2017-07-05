package main

import (
	"net/http"
	"log"
	"time"
	"github.com/dchest/uniuri"
	"sync"
	"path/filepath"
	"golang.org/x/crypto/scrypt"
	"encoding/base64"
)

var unAuthenticatedURLS = []string{"favicon.ico", "/login", "*robots.txt"}

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

func createSessionCookie(w http.ResponseWriter) *sessionData{

	sesID := uniuri.NewLen(32)

	sessionCookie := http.Cookie{Name: "session", Value: sesID, Path: "/"}

	data := NewSessionData()

	ssMutex.Lock()
	sessions[sesID] = data
	ssMutex.Unlock()


	http.SetCookie(w, &sessionCookie)

	return data
}


func cleanSessionStore(){

	if lastSessionStoreClean.Add(config.SessionStoreCleanInterval).After(time.Now()) {
		return
	}

	lastSessionStoreClean = time.Now()

	ssMutex.Lock()

	log.Println("cleaning session store...")
	log.Println(sessions)

	for k, v := range sessions {
		if v.lastUsed.Add(config.SessionExpirationTime).Before(time.Now()) {
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
		if val.lastUsed.Add(config.SessionExpirationTime).After(time.Now()) {

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
	redirectAfterLogin := http.Cookie{Name: "afterLogin", Value: r.URL.Path, Path: "/"}
	http.SetCookie(w, &redirectAfterLogin)

	log.Print("redirecting to: "+"/login ")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
					// expire the cookie immediately
					redirect.Expires = time.Now().Add(- time.Hour*24)
					http.SetCookie(w, redirect)

					http.Redirect(w, r, redirect.Value, http.StatusSeeOther)
					log.Printf("redirecting to: %s ", redirect.Value)
				}else{
					http.Redirect(w, r, config.HomePage, http.StatusSeeOther)
					log.Printf("redirecting to %s", config.HomePage)
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
			log.Print(err)
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
				// expire the cookie immediately
				redirect.Expires = time.Now().Add(- time.Hour*24)
				http.SetCookie(w, redirect)

				http.Redirect(w, r, redirect.Value, http.StatusSeeOther)
				log.Printf("redirecting to: %s ", redirect.Value)
			}else{
				http.Redirect(w, r, config.HomePage, http.StatusSeeOther)
				log.Printf("redirecting to %s", config.HomePage)
			}

		} else if (user == "" && pass == "") {
			log.Print("login POST with empty user and pass. Assuming a derp. Redirecting to /login with 303")
			http.Redirect(w, r, "/login", http.StatusSeeOther);
		} else {
			log.Print("Wrong username or password")
			http.Redirect(w, r, "/login", http.StatusSeeOther);
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	getSessionCookie(w, r).authenticated = false
	// TODO: Delete token as well

	http.ServeFile(w, r, "html/logout.html")
}

func matchesAny(str string, pattern []string) bool {

	for _, patrn := range pattern {

		if match, err := filepath.Match(patrn, str); match {
			if (err != nil) {
				log.Print(err)
				return false
			}

			return true
		}

	}
	return false
}

//Authentication Middleware
func authMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("authRequest to %s", r.URL.Path)

		cookie := getSessionCookie(w,r)

		if (matchesAny(r.URL.Path, unAuthenticatedURLS)) {
			log.Println("Request URL in the allowed requests list. Passing authMW.")
			next.ServeHTTP(w, r)
			return;
		}

		if(cookie.authenticated){
			log.Println("Authenticated")
			next.ServeHTTP(w, r)
		}else{
			log.Println("Unauthenticated")
			redirectLogin(w, r)
		}

	})
}

func hashPassword(password, salt []byte) string {
	val, err := scrypt.Key(password, salt, 32768, 8, 1, 32)
	if err != nil {
		panic(err)
	}

	str := base64.StdEncoding.EncodeToString([]byte(val))
	//str := hex.EncodeToString(val)

	return str
}
