package handlers

import (
	"net/http"
	"starter-go-app/templates"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Header.Get("session")

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		templates.Tmpl.ExecuteTemplate(w, "404.html", map[string]interface{}{
			"Session": session,
			"Title":   "Page Not Found",
		})
		return
	}

	templates.Tmpl.ExecuteTemplate(w, "home.html", map[string]interface{}{
		"Session": session,
		"Title":   "Home",
	})
}
