package main

import (
	"log"
	"net/http"
	"os"
	_ "starter-go-app/database"
	"starter-go-app/handlers"

	"github.com/go-http-utils/cookie"
)

func middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add session to context
		cookies := cookie.New(w, r, os.Getenv("COOKIE_KEY"))
		session, _ := cookies.Get("session", true)

		r.Header.Set("session", session)

		// Log request
		log.Printf("Request: %v %v\n", r.Method, r.URL.Path)

		// Serve Next HTTP
		handler.ServeHTTP(w, r)
	})
}

func auth(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := cookie.New(w, r, os.Getenv("COOKIE_KEY"))
		session, err := cookies.Get("session", true)

		if err != nil {
			http.Redirect(w, r, "/users/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}

		if session == "" {
			http.Redirect(w, r, "/users/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}

		r.Header.Set("session", session)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	r := http.NewServeMux()

	// Static assets
	r.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Pages
	r.HandleFunc("/", handlers.HomeHandler)

	// Users
	r.HandleFunc("/users/login", handlers.UserLoginHandler)
	r.HandleFunc("/users/register", handlers.UserRegisterHandler)
	r.HandleFunc("/users/verify", handlers.UserVerifyHandler)
	r.HandleFunc("/users/logout", handlers.UserLogoutHandler)
	r.HandleFunc("/users/forgot", handlers.UserForgotHandler)
	r.HandleFunc("/users/reset", handlers.UserResetHandler)

	// Posts
	r.HandleFunc("/posts", auth(handlers.PostsIndexHandler))
	r.HandleFunc("/posts/new", auth(handlers.PostsNewHandler))
	r.HandleFunc("/posts/view/{id}", auth(handlers.PostsViewHandler))
	r.HandleFunc("/posts/delete/{id}", auth(handlers.PostsDeleteHandler))
	r.HandleFunc("/posts/edit/{id}", auth(handlers.PostsEditHandler))

	http.ListenAndServe(":8080", middleware(r))
}
