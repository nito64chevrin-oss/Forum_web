package views

import (
	"database/sql"
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// ========================================
// STRUCTURES
// ========================================

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

// ========================================
// HANDLER : INSCRIPTION
// ========================================

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

		// Générer un UUID pour l'utilisateur
		userID := uuid.NewV4().String()

		// Insérer dans la base de données
		query := `
			INSERT INTO users (id, username, email, password, favorite_jojo_part, favorite_stand, created_at)
			VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		`

		_, err = db.Exec(query,
			userID,
			req.Username,
			req.Email,
			string(hashedPassword),
			req.FavoriteJojoPart,
			req.FavoriteStand,
		)

		if err != nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email ou nom d'utilisateur déjà utilisé",
			})
			return
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
				"username": req.Username,
				"email":    req.Email,
			},
		})
	}
}

// ========================================
// HANDLER : CONNEXION
// ========================================

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Décoder les données JSON
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Données invalides",
			})
			return
		}

		// Récupérer l'utilisateur depuis la base de données
		var userID, username, hashedPassword string
		query := "SELECT id, username, password FROM users WHERE email = ?"

		err := db.QueryRow(query, req.Email).Scan(&userID, &username, &hashedPassword)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email ou mot de passe incorrect",
			})
			return
		}

		// Vérifier le mot de passe avec Bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email ou mot de passe incorrect",
			})
			return
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

// ========================================
// HANDLER : RÉCUPÉRER LE PROFIL UTILISATEUR
// ========================================

func GetUserProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Récupérer l'email depuis l'URL (/api/user?email=...)
		email := r.URL.Query().Get("email")

		if email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email manquant",
			})
			return
		}

		// Structure pour stocker les données utilisateur
		var user struct {
			ID               string `json:"id"`
			Username         string `json:"username"`
			Email            string `json:"email"`
			FavoriteJojoPart string `json:"favorite_jojo_part"`
			FavoriteStand    string `json:"favorite_stand"`
			Bio              string `json:"bio"`
			Location         string `json:"location"`
			AvatarURL        string `json:"avatar_url"`
			CreatedAt        string `json:"created_at"`
		}

		// Requête SQL pour récupérer le profil
		query := `
			SELECT id, username, email, 
			       COALESCE(favorite_jojo_part, '') as favorite_jojo_part,
			       COALESCE(favorite_stand, '') as favorite_stand,
			       COALESCE(bio, '') as bio,
			       COALESCE(location, '') as location,
			       COALESCE(avatar_url, '') as avatar_url,
			       created_at
			FROM users 
			WHERE email = ?
		`

		err := db.QueryRow(query, email).Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FavoriteJojoPart,
			&user.FavoriteStand,
			&user.Bio,
			&user.Location,
			&user.AvatarURL,
			&user.CreatedAt,
		)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Utilisateur non trouvé",
			})
			return
		}

		// Calculer les stats
		var stats struct {
			PostsCount    int `json:"posts_count"`
			CommentsCount int `json:"comments_count"`
			LikesReceived int `json:"likes_received"`
		}

		// Compter les posts
		db.QueryRow(`SELECT COUNT(*) FROM posts WHERE user_id = ?`, user.ID).Scan(&stats.PostsCount)

		// Compter les commentaires
		db.QueryRow(`SELECT COUNT(*) FROM comments WHERE user_id = ?`, user.ID).Scan(&stats.CommentsCount)

		// Compter les likes reçus (sur les posts + sur les commentaires)
		db.QueryRow(`
			SELECT COUNT(*) FROM reactions r
			JOIN posts p ON r.post_id = p.id
			WHERE p.user_id = ? AND r.reaction_type = 'like'
		`, user.ID).Scan(&stats.LikesReceived)

		// Renvoyer les données en JSON
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"user":   user,
			"stats":  stats,
		})
	}
}

// ========================================
// HANDLER : METTRE À JOUR LE PROFIL
// ========================================

func UpdateUserProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Décoder les données JSON
		var req struct {
			Email            string `json:"email"`
			Username         string `json:"username"`
			Bio              string `json:"bio"`
			Location         string `json:"location"`
			AvatarURL        string `json:"avatar_url"`
			FavoriteStand    string `json:"favorite_stand"`
			FavoriteJojoPart string `json:"favorite_jojo_part"`
			Interests        string `json:"interests"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Données invalides",
			})
			return
		}

		// Vérifier que l'email existe
		if req.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email requis",
			})
			return
		}

		// Construire la requête UPDATE
		query := `
			UPDATE users 
			SET 
				username = COALESCE(?, username),
				bio = COALESCE(?, bio),
				location = COALESCE(?, location),
				avatar_url = COALESCE(?, avatar_url),
				favorite_stand = COALESCE(?, favorite_stand),
				favorite_jojo_part = COALESCE(?, favorite_jojo_part),
				interests = COALESCE(?, interests)
			WHERE email = ?
		`

		result, err := db.Exec(query,
			nullIfEmpty(req.Username),
			nullIfEmpty(req.Bio),
			nullIfEmpty(req.Location),
			nullIfEmpty(req.AvatarURL),
			nullIfEmpty(req.FavoriteStand),
			nullIfEmpty(req.FavoriteJojoPart),
			nullIfEmpty(req.Interests),
			req.Email,
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erreur serveur",
			})
			return
		}

		// Vérifier que la mise à jour a eu lieu
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Utilisateur non trouvé",
			})
			return
		}

		// Récupérer les données mises à jour
		var user struct {
			ID               string `json:"id"`
			Username         string `json:"username"`
			Email            string `json:"email"`
			FavoriteJojoPart string `json:"favorite_jojo_part"`
			FavoriteStand    string `json:"favorite_stand"`
			Bio              string `json:"bio"`
			Location         string `json:"location"`
			AvatarURL        string `json:"avatar_url"`
			CreatedAt        string `json:"created_at"`
		}

		getQuery := `
			SELECT id, username, email, 
			       COALESCE(favorite_jojo_part, '') as favorite_jojo_part,
			       COALESCE(favorite_stand, '') as favorite_stand,
			       COALESCE(bio, '') as bio,
			       COALESCE(location, '') as location,
			       COALESCE(avatar_url, '') as avatar_url,
			       created_at
			FROM users 
			WHERE email = ?
		`

		err = db.QueryRow(getQuery, req.Email).Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FavoriteJojoPart,
			&user.FavoriteStand,
			&user.Bio,
			&user.Location,
			&user.AvatarURL,
			&user.CreatedAt,
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erreur serveur",
			})
			return
		}

		// Renvoyer les données mises à jour
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"user":   user,
		})
	}
}

// Fonction helper pour retourner NULL si vide
func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
