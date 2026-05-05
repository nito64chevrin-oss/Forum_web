package views

import (
	"database/sql"
	"encoding/json"
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

		// Décoder les données JSON
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Données invalides",
			})
			return
		}

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

		var user struct {
			ID               string   `json:"id"`
			Username         string   `json:"username"`
			Email            string   `json:"email"`
			FavoriteJojoPart string   `json:"favorite_jojo_part"`
			FavoriteStand    string   `json:"favorite_stand"`
			Bio              string   `json:"bio"`
			Location         string   `json:"location"`
			AvatarURL        string   `json:"avatar_url"`
			CreatedAt        string   `json:"created_at"`
			Interests        []string `json:"interests"`
		}

		query := `
			SELECT id, username, email, 
			       COALESCE(favorite_jojo_part, '') as favorite_jojo_part,
			       COALESCE(favorite_stand, '') as favorite_stand,
			       COALESCE(bio, '') as bio,
			       COALESCE(location, '') as location,
			       COALESCE(avatar_url, '') as avatar_url,
			       COALESCE(interests, '[]') as interests,
			       created_at
			FROM users 
			WHERE email = ?
		`

		var interestsJSON string

		err := db.QueryRow(query, email).Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FavoriteJojoPart,
			&user.FavoriteStand,
			&user.Bio,
			&user.Location,
			&user.AvatarURL,
			&interestsJSON,
			&user.CreatedAt,
		)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Utilisateur non trouvé",
			})
			return
		}

		json.Unmarshal([]byte(interestsJSON), &user.Interests)

		// Renvoyer les données
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"user":   user,
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
