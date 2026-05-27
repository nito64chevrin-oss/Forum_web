package main

import (
	"database/sql"
	"forum/views"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Ouvrir la base de données
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Routes AUTH
	http.HandleFunc("/api/register", views.RegisterHandler(db))
	http.HandleFunc("/api/login", views.LoginHandler(db))
	http.HandleFunc("/api/user", views.GetUserProfileHandler(db))
	http.HandleFunc("/api/user/update", views.UpdateUserProfileHandler(db))

	// Routes POSTS
	http.HandleFunc("/api/posts", views.CreatePostHandler(db))
	http.HandleFunc("/api/posts/list", views.GetPostsHandler(db))
	http.HandleFunc("/api/posts/get", views.GetPostByIDHandler(db))

	// Fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))

	// Pages HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			views.CreateCommentHandler(db)(w, r)
		} else {
			views.GetCommentsHandler(db)(w, r)
		}
	})

	http.HandleFunc("/category.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/category.html")
	})

	http.HandleFunc("/category.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "static/category.css")
	})

	http.HandleFunc("/category.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "static/category.js")
	})

	// Like commentaire
	http.HandleFunc("/api/comments/like", views.LikeCommentHandler(db))

	// Favori commentaire
	http.HandleFunc("/api/favorites", views.ToggleFavoriteHandler(db))

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "static/style.css")
	})

	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "static/script.js")
	})

	http.HandleFunc("/auth.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/auth.html")
	})

	http.HandleFunc("/auth.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "static/auth.css")
	})

	http.HandleFunc("/auth.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "static/auth.js")
	})

	http.HandleFunc("/profile.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/profile.html")
	})

	http.HandleFunc("/profile.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "static/profile.css")
	})

	http.HandleFunc("/profile.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "static/profile.js")
	})

	// NOUVEAU - Pages pour créer un post
	http.HandleFunc("/create-post.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/create-post.html")
	})

	http.HandleFunc("/create-post.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "static/create-post.css")
	})

	http.HandleFunc("/create-post.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "static/create-post.js")
	})

	http.HandleFunc("/posts-display.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "static/posts-display.css")
	})

	// Routes fichiers statiques
	http.HandleFunc("/post.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/post.html")
	})

	http.HandleFunc("/post.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "static/post.css")
	})

	http.HandleFunc("/post.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "static/post.js")
	})

	http.HandleFunc("/posts-loader.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "static/posts-loader.js")
	})

	// Lancer le serveur
	log.Println(" Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
