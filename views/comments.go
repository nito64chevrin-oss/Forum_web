package views

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// ========================================
// STRUCTURES
// ========================================

type CreateCommentRequest struct {
	PostID      string `json:"post_id"`
	AuthorEmail string `json:"author_email"`
	Content     string `json:"content"`
	ImageURL    string `json:"image_url"`
}

type Comment struct {
	ID         string `json:"id"`
	PostID     string `json:"post_id"`
	AuthorID   string `json:"author_id"`
	AuthorName string `json:"author_name"`
	AvatarURL  string `json:"avatar_url"`
	Content    string `json:"content"`
	ImageURL   string `json:"image_url"`
	CreatedAt  string `json:"created_at"`
	LikesCount int    `json:"likes_count"`
	IsFavorite bool   `json:"is_favorite"`
}

// ========================================
// HANDLER : CRÉER UN COMMENTAIRE
// ========================================

func CreateCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Méthode non autorisée",
			})
			return
		}

		// Décoder les données JSON
		var req CreateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("❌ Erreur décodage JSON:", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Données invalides",
			})
			return
		}

		log.Println("💬 Création commentaire sur post:", req.PostID)

		// Récupérer l'ID de l'auteur depuis son email
		var userID string
		err := db.QueryRow("SELECT id FROM users WHERE email = ?", req.AuthorEmail).Scan(&userID)
		if err != nil {
			log.Println("❌ Auteur non trouvé:", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Utilisateur non connecté",
			})
			return
		}

		// Générer un UUID pour le commentaire
		commentID := uuid.NewV4().String()

		// Insérer le commentaire dans la base de données
		query := `
			INSERT INTO comments (id, post_id, user_id, content, image_url, created_at)
			VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		`

		_, err = db.Exec(query, commentID, req.PostID, userID, req.Content, req.ImageURL)
		if err != nil {
			log.Println("❌ Erreur insertion commentaire:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erreur lors de la création du commentaire",
			})
			return
		}

		log.Println("✅ Commentaire créé avec succès:", commentID)

		// Réponse succès
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "success",
			"comment_id": commentID,
		})
	}
}

// ========================================
// HANDLER : RÉCUPÉRER LES COMMENTAIRES D'UN POST
// ========================================

func GetCommentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Paramètres de requête
		postID := r.URL.Query().Get("post_id")
		userEmail := r.URL.Query().Get("user_email") // Pour vérifier les favoris

		if postID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "post_id requis",
			})
			return
		}

		log.Println("💬 Récupération commentaires pour post:", postID)

		// Récupérer l'ID utilisateur s'il est connecté
		var currentUserID string
		if userEmail != "" {
			db.QueryRow("SELECT id FROM users WHERE email = ?", userEmail).Scan(&currentUserID)
		}

		// Requête SQL avec comptage des likes
		query := `
			SELECT 
				c.id,
				c.post_id,
				c.user_id,
				u.username as author_name,
				COALESCE(u.avatar_url, 'static/images/default-avatar.png') as avatar_url,
				c.content,
				COALESCE(c.image_url, '') as image_url,
				c.created_at,
				COALESCE(COUNT(DISTINCT cr.id), 0) as likes_count,
				CASE WHEN f.id IS NOT NULL THEN 1 ELSE 0 END as is_favorite
			FROM comments c
			INNER JOIN users u ON c.user_id = u.id
			LEFT JOIN comment_reactions cr ON c.id = cr.comment_id AND cr.reaction_type = 'like'
			LEFT JOIN favorites f ON c.id = f.comment_id AND f.user_id = ?
			WHERE c.post_id = ?
			GROUP BY c.id, c.post_id, c.user_id, u.username, u.avatar_url, c.content, c.created_at, f.id
			ORDER BY c.created_at ASC
		`

		rows, err := db.Query(query, currentUserID, postID)
		if err != nil {
			log.Println("❌ Erreur récupération commentaires:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erreur lors de la récupération des commentaires",
			})
			return
		}
		defer rows.Close()

		// Parcourir les résultats
		var comments []Comment

		for rows.Next() {
			var comment Comment
			err := rows.Scan(
				&comment.ID,
				&comment.PostID,
				&comment.AuthorID,
				&comment.AuthorName,
				&comment.AvatarURL,
				&comment.Content,
				&comment.ImageURL,
				&comment.CreatedAt,
				&comment.LikesCount,
				&comment.IsFavorite,
			)

			if err != nil {
				log.Println("❌ Erreur scan commentaire:", err)
				continue
			}

			comments = append(comments, comment)
		}

		log.Println("✅ Commentaires récupérés:", len(comments))

		// Réponse
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "success",
			"comments": comments,
		})
	}
}

// ========================================
// HANDLER : LIKER UN COMMENTAIRE
// ========================================

func LikeCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "POST requis"})
			return
		}

		var req struct {
			CommentID   string `json:"comment_id"`
			AuthorEmail string `json:"author_email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Données invalides"})
			return
		}

		// Récupérer l'ID utilisateur
		var userID string
		err := db.QueryRow("SELECT id FROM users WHERE email = ?", req.AuthorEmail).Scan(&userID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Utilisateur non connecté"})
			return
		}

		// Vérifier si le like existe déjà
		var existingID string
		err = db.QueryRow("SELECT id FROM comment_reactions WHERE comment_id = ? AND user_id = ?", req.CommentID, userID).Scan(&existingID)

		if err == nil {
			// Le like existe, le supprimer
			_, err := db.Exec("DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ?", req.CommentID, userID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Erreur suppression like"})
				return
			}
			log.Println("❌ Like supprimé pour commentaire:", req.CommentID)
			json.NewEncoder(w).Encode(map[string]string{"status": "unliked"})
		} else {
			// Le like n'existe pas, l'ajouter
			reactionID := uuid.NewV4().String()
			_, err := db.Exec(`
				INSERT INTO comment_reactions (id, comment_id, user_id, reaction_type, created_at)
				VALUES (?, ?, ?, 'like', CURRENT_TIMESTAMP)
			`, reactionID, req.CommentID, userID)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Erreur ajout like"})
				return
			}
			log.Println("❤️ Like ajouté pour commentaire:", req.CommentID)
			json.NewEncoder(w).Encode(map[string]string{"status": "liked"})
		}
	}
}

// ========================================
// HANDLER : AJOUTER/SUPPRIMER UN FAVORI
// ========================================

func ToggleFavoriteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "POST requis"})
			return
		}

		var req struct {
			CommentID   string `json:"comment_id"`
			AuthorEmail string `json:"author_email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Données invalides"})
			return
		}

		// Récupérer l'ID utilisateur
		var userID string
		err := db.QueryRow("SELECT id FROM users WHERE email = ?", req.AuthorEmail).Scan(&userID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Utilisateur non connecté"})
			return
		}

		// Vérifier si le favori existe
		var existingID string
		err = db.QueryRow("SELECT id FROM favorites WHERE comment_id = ? AND user_id = ?", req.CommentID, userID).Scan(&existingID)

		if err == nil {
			// Le favori existe, le supprimer
			_, err := db.Exec("DELETE FROM favorites WHERE comment_id = ? AND user_id = ?", req.CommentID, userID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Erreur suppression favori"})
				return
			}
			log.Println("🗑️ Favori supprimé pour commentaire:", req.CommentID)
			json.NewEncoder(w).Encode(map[string]string{"status": "unfavorited"})
		} else {
			// Le favori n'existe pas, l'ajouter
			favoriteID := uuid.NewV4().String()
			_, err := db.Exec(`
				INSERT INTO favorites (id, comment_id, user_id, created_at)
				VALUES (?, ?, ?, CURRENT_TIMESTAMP)
			`, favoriteID, req.CommentID, userID)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Erreur ajout favori"})
				return
			}
			log.Println("⭐ Favori ajouté pour commentaire:", req.CommentID)
			json.NewEncoder(w).Encode(map[string]string{"status": "favorited"})
		}
	}
}
