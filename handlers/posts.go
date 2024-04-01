package handlers

import (
	"net/http"
	"starter-go-app/database"
	"starter-go-app/templates"
)

func PostsIndexHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Header.Get("session")
	posts := []database.Posts{}

	// sort newest to oldest
	database.DB.Order("id desc").Find(&posts)

	templates.Tmpl.ExecuteTemplate(w, "posts.html", map[string]interface{}{
		"Session": session,
		"Title":   "Posts Index",
		"Posts":   posts,
	})
}

func PostsViewHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	session := r.Header.Get("session")
	post := database.Posts{}

	database.DB.Find(&post, "id = ?", id)

	templates.Tmpl.ExecuteTemplate(w, "posts-view.html", map[string]interface{}{
		"Session": session,
		"ID":      id,
		"Title":   post.Title,
		"Body":    post.Body,
	})
}

func PostsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	database.DB.Delete(&database.Posts{}, "id = ?", id)

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func PostsEditHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	session := r.Header.Get("session")

	if r.Method == "POST" {
		title := r.FormValue("title")
		body := r.FormValue("body")

		if title == "" || body == "" {
			http.Redirect(w, r, "/posts", http.StatusSeeOther)
			return
		}

		database.DB.Save(&database.Posts{ID: id, Title: title, Body: body})

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	} else {
		post := database.Posts{}

		database.DB.Find(&post, "id = ?", id)

		templates.Tmpl.ExecuteTemplate(w, "posts-edit.html", map[string]interface{}{
			"Session": session,
			"Title":   post.Title,
			"Body":    post.Body,
		})
	}
}

func PostsNewHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Header.Get("session")

	if r.Method == "POST" {
		title := r.FormValue("title")
		body := r.FormValue("body")

		if title == "" || body == "" {
			http.Redirect(w, r, "/posts", http.StatusSeeOther)
			return
		}

		database.DB.Create(&database.Posts{Title: title, Body: body})

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}

	templates.Tmpl.ExecuteTemplate(w, "posts-new.html", map[string]interface{}{
		"Session": session,
		"Title":   "New Post",
	})
}
