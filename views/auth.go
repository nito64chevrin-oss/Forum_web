package views

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username         string   `json:"username"`
	Email            string   `json:"email"`
	Password         string   `json:"password"`
	FavoriteJojoPart string   `json:"favorite_jojo_part"`
	FavoriteStand    string   `json:"favorite_stand"`
	Interests        []string `json:"interests"`
	Avatar           string   `json:"avatar"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Décoder les données JSON
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Données invalides",
			})
			return
		}

		// Hasher le mot de passe avec Bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erreur serveur",
			})
			return
		}
		userID := uuid.NewV4().String()

		interestsJSON, _ := json.Marshal(req.Interests)

		query := `
			INSERT INTO users (id, username, email, password, favorite_jojo_part, favorite_stand, interests, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		`

		_, err = db.Exec(query,
			userID,
			req.Username,
			req.Email,
			string(hashedPassword),
			req.FavoriteJojoPart,
			req.FavoriteStand,
			string(interestsJSON),
		)

		if err != nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email ou nom d'utilisateur déjà utilisé",
			})
			return
		}

		sessionID := uuid.NewV4().String()
		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  sessionID,
			MaxAge: 2592000, // 30 jours
			Path:   "/",
		})

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"user": map[string]string{
				"id":       userID,
				"username": req.Username,
				"email":    req.Email,
			},
		})
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Println("🔵 LoginHandler appelé") // ✅ LOG 1

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("❌ Erreur décodage JSON:", err) // ✅ LOG 2
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Données invalides",
			})
			return
		}

		log.Println("📧 Email reçu:", req.Email)       // ✅ LOG 3
		log.Println("🔑 Password reçu:", req.Password) // ✅ LOG 4

		var userID, username, hashedPassword string
		query := "SELECT id, username, password FROM users WHERE email = ?"

		err := db.QueryRow(query, req.Email).Scan(&userID, &username, &hashedPassword)

		if err != nil {
			log.Println("❌ User non trouvé:", err) // ✅ LOG 5
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email ou mot de passe incorrect",
			})
			return
		}

		log.Println("✅ User trouvé:", username)                    // ✅ LOG 6
		log.Println("🔐 Hash de la DB:", hashedPassword[:20]+"...") // ✅ LOG 7

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
		if err != nil {
			log.Println("❌ Mot de passe FAUX pour:", req.Email)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email ou mot de passe incorrect",
			})
			return
		} else {
			log.Println("✅ Mot de passe CORRECT pour:", req.Email)
		}

		// Créer un cookie de session
		sessionID := uuid.NewV4().String()
		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  sessionID,
			MaxAge: 2592000, // 30 jours
			Path:   "/",
		})

		// Réponse succès
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"user": map[string]string{
				"id":       userID,
				"username": username,
				"email":    req.Email,
			},
		})
	}
}

func GetUserProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Récupérer l'email depuis les paramètres de requête
		email := r.URL.Query().Get("email")

		if email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email requis",
			})
			return
		}

		// Requête pour récupérer l'utilisateur
		var user struct {
			ID        string `json:"id"`
			Username  string `json:"username"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
			Bio       string `json:"bio"`
			CreatedAt string `json:"created_at"`
		}

		query := `
			SELECT id, username, email, 
			       COALESCE(avatar_url, 'static/images/default-avatar.png') as avatar_url,
			       COALESCE(bio, '') as bio,
			       created_at
			FROM users
			WHERE email = ?
		`

		err := db.QueryRow(query, email).Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.AvatarURL,
			&user.Bio,
			&user.CreatedAt,
		)

		if err != nil {
			log.Println("❌ Utilisateur non trouvé:", err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Utilisateur non trouvé",
			})
			return
		}

		// Compter le nombre de posts
		var postsCount int
		db.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", user.ID).Scan(&postsCount)

		// Compter le nombre de commentaires
		var commentsCount int
		db.QueryRow("SELECT COUNT(*) FROM comments WHERE user_id = ?", user.ID).Scan(&commentsCount)

		// Compter le nombre de likes reçus sur les posts
		var likesReceivedOnPosts int
		db.QueryRow(`
			SELECT COUNT(*) 
			FROM reactions r
			INNER JOIN posts p ON r.post_id = p.id
			WHERE p.user_id = ? AND r.reaction_type = 'like'
		`, user.ID).Scan(&likesReceivedOnPosts)

		// Compter le nombre de likes reçus sur les commentaires
		var likesReceivedOnComments int
		db.QueryRow(`
			SELECT COUNT(*)
			FROM reactions r
			INNER JOIN comments c ON r.comment_id = c.id
			WHERE c.user_id = ? AND r.reaction_type = 'like'
		`, user.ID).Scan(&likesReceivedOnComments)

		// Total des likes reçus
		totalLikesReceived := likesReceivedOnPosts + likesReceivedOnComments

		log.Println("✅ Profil récupéré:", user.Username, "- Posts:", postsCount, "- Commentaires:", commentsCount, "- Likes:", totalLikesReceived)

		// Réponse avec stats
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"user":   user,
			"stats": map[string]int{
				"posts_count":       postsCount,
				"comments_count":    commentsCount,
				"likes_received":    totalLikesReceived,
				"likes_on_posts":    likesReceivedOnPosts,
				"likes_on_comments": likesReceivedOnComments,
			},
		})
	}
}

func UpdateUserProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Décoder les données JSON
		var req struct {
			Email            string `json:"email"`
			Bio              string `json:"bio"`
			Location         string `json:"location"`
			FavoriteJojoPart string `json:"favorite_jojo_part"`
			FavoriteStand    string `json:"favorite_stand"`
			AvatarURL        string `json:"avatar_url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Données invalides",
			})
			return
		}

		// Mettre à jour dans SQL
		query := `
			UPDATE users 
			SET bio = ?, location = ?, favorite_jojo_part = ?, favorite_stand = ?, avatar_url = ?
			WHERE email = ?
		`

		_, err := db.Exec(query,
			req.Bio,
			req.Location,
			req.FavoriteJojoPart,
			req.FavoriteStand,
			req.AvatarURL,
			req.Email,
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erreur lors de la mise à jour",
			})
			return
		}

		// Réponse succès
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
		})
	}
}
