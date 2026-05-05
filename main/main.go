package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum/views"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "./static/style.css")
	})
	http.HandleFunc("/auth.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "./static/auth.css")
	})
	http.HandleFunc("/profile.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "./static/profile.css")
	})

	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "./static/script.js")
	})
	http.HandleFunc("/auth.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "./static/auth.js")
	})
	http.HandleFunc("/profile.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "./static/profile.js")
	})

	http.HandleFunc("/api/register", views.RegisterHandler(db))
	http.HandleFunc("/api/login", views.LoginHandler(db))
	http.HandleFunc("/api/user", views.GetUserProfileHandler(db))
	http.HandleFunc("/api/user/update", views.UpdateUserProfileHandler(db))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./static/images"))))

	http.HandleFunc("/auth.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/auth.html")
	})
	http.HandleFunc("/profile.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/profile.html")
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./static/index.html")
	})

	log.Println("🚀 Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
