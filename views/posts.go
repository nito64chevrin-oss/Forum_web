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

type CreatePostRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	CategoryID  string   `json:"category_id"`  // Catégorie principale
	CategoryIDs []string `json:"category_ids"` // Toutes les catégories sélectionnées
	AuthorEmail string   `json:"author_email"`
	ImageURL    string   `json:"image_url"`
	Tags        []string `json:"tags"`
}

type Post struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	CategoryID    string   `json:"category_id"`
	AuthorID      string   `json:"author_id"`
	AuthorName    string   `json:"author_name"`
	AvatarURL     string   `json:"avatar_url"`
	ImageURL      string   `json:"image_url"`
	CreatedAt     string   `json:"created_at"`
	LikesCount    int      `json:"likes_count"`
	CommentsCount int      `json:"comments_count"`
	Tags          []string `json:"tags"`
}

// ========================================
// HANDLER : CRÉER UN POST
// ========================================

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Décoder les données JSON
		var req CreatePostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("❌ Erreur décodage JSON:", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Données invalides",
			})
			return
		}

		log.Println("📝 Création post par:", req.AuthorEmail)

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

		// Générer un UUID pour le post
		postID := uuid.NewV4().String()

		// Insérer le post dans la base de données
		query := `
			INSERT INTO posts (id, title, content, category_id, user_id, image_url, created_at)
			VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		`

		_, err = db.Exec(query, postID, req.Title, req.Content, req.CategoryID, userID, req.ImageURL)
		if err != nil {
			log.Println("❌ Erreur insertion post:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erreur lors de la création du post",
			})
			return
		}

		// Ajouter les catégories multiples dans post_categories
		if len(req.CategoryIDs) > 0 {
			for _, catID := range req.CategoryIDs {
				_, err := db.Exec("INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, catID)
				if err != nil {
					log.Println("❌ Erreur liaison post-catégorie:", err)
				}
			}
			log.Println("✅ Post ajouté dans", len(req.CategoryIDs), "catégories")
		}

		// Ajouter les tags
		if len(req.Tags) > 0 {
			for _, tagName := range req.Tags {
				// Vérifier si le tag existe
				var tagID string
				err := db.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)

				if err == sql.ErrNoRows {
					// Le tag n'existe pas, on le crée avec un UUID
					tagID = uuid.NewV4().String()
					_, err := db.Exec("INSERT INTO tags (id, name) VALUES (?, ?)", tagID, tagName)
					if err != nil {
						log.Println("❌ Erreur création tag:", err)
						continue
					}
				} else if err != nil {
					log.Println("❌ Erreur recherche tag:", err)
					continue
				}

				// Lier le tag au post
				_, err = db.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", postID, tagID)
				if err != nil {
					log.Println("❌ Erreur liaison post-tag:", err)
				}
			}
		}

		log.Println("✅ Post créé avec succès:", postID)

		// Réponse succès
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"post_id": postID,
		})
	}
}

// ========================================
// HANDLER : RÉCUPÉRER LES POSTS
// ========================================

func GetPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Paramètres de requête
		categoryID := r.URL.Query().Get("category_id")
		sortBy := r.URL.Query().Get("sort") // "recent" ou "popular"

		if sortBy == "" {
			sortBy = "recent"
		}

		log.Println("📋 Récupération posts - Catégorie:", categoryID, "- Tri:", sortBy)

		// Construire la requête SQL
		query := `
			SELECT 
				p.id, 
				p.title, 
				p.content, 
				p.category_id, 
				p.user_id,
				u.username as author_name,
				COALESCE(u.avatar_url, 'static/images/default-avatar.png') as avatar_url,
				COALESCE(p.image_url, '') as image_url,
				p.created_at,
				COALESCE(COUNT(DISTINCT r.id), 0) as likes_count,
				COALESCE(COUNT(DISTINCT c.id), 0) as comments_count
			FROM posts p
			INNER JOIN users u ON p.user_id = u.id
			LEFT JOIN reactions r ON r.post_id = p.id AND r.reaction_type = 'like'
			LEFT JOIN comments c ON c.post_id = p.id
		`

		// Filtrer par catégorie si spécifié (supporte les deux systèmes)
		if categoryID != "" {
			query += ` 
				WHERE (
					p.category_id = ? 
					OR EXISTS (
						SELECT 1 FROM post_categories pc 
						WHERE pc.post_id = p.id AND pc.category_id = ?
					)
				)
			`
		}

		query += " GROUP BY p.id, p.title, p.content, p.category_id, p.user_id, u.username, u.avatar_url, p.image_url, p.created_at"

		// Trier
		if sortBy == "popular" {
			query += " ORDER BY likes_count DESC, created_at DESC"
		} else {
			query += " ORDER BY p.created_at DESC"
		}

		// Exécuter la requête
		var rows *sql.Rows
		var err error

		if categoryID != "" {
			rows, err = db.Query(query, categoryID, categoryID) // ✅ Deux fois pour les deux ?
		} else {
			rows, err = db.Query(query)
		}

		if err != nil {
			log.Println("❌ Erreur récupération posts:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erreur serveur",
			})
			return
		}
		defer rows.Close()

		// Parser les résultats
		posts := []Post{}
		for rows.Next() {
			var post Post
			err := rows.Scan(
				&post.ID,
				&post.Title,
				&post.Content,
				&post.CategoryID,
				&post.AuthorID,
				&post.AuthorName,
				&post.AvatarURL,
				&post.ImageURL,
				&post.CreatedAt,
				&post.LikesCount,
				&post.CommentsCount,
			)
			if err != nil {
				log.Println("❌ Erreur scan post:", err)
				continue
			}

			// Récupérer les tags du post
			tagRows, err := db.Query(`
				SELECT t.name 
				FROM tags t
				JOIN post_tags pt ON t.id = pt.tag_id
				WHERE pt.post_id = ?
			`, post.ID)

			if err == nil {
				tags := []string{}
				for tagRows.Next() {
					var tag string
					tagRows.Scan(&tag)
					tags = append(tags, tag)
				}
				post.Tags = tags
				tagRows.Close()
			}

			posts = append(posts, post)
		}

		log.Println("✅ Posts récupérés:", len(posts))

		// Renvoyer les posts
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"posts":  posts,
		})
	}
}

// ========================================
// HANDLER : RÉCUPÉRER UN POST PAR ID
// ========================================

func GetPostByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Récupérer l'ID depuis l'URL
		postID := r.URL.Query().Get("id")

		if postID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "ID manquant",
			})
			return
		}

		log.Println("📖 Récupération post:", postID)

		// Requête SQL
		query := `
			SELECT 
				p.id, 
				p.title, 
				p.content, 
				p.category_id, 
				p.user_id,
				u.username as author_name,
				COALESCE(u.avatar_url, 'static/images/default-avatar.png') as avatar_url,
				COALESCE(p.image_url, '') as image_url,
				p.created_at,
				COALESCE(COUNT(DISTINCT r.id), 0) as likes_count,
				COALESCE(COUNT(DISTINCT c.id), 0) as comments_count
			FROM posts p
			INNER JOIN users u ON p.user_id = u.id
			LEFT JOIN reactions r ON r.post_id = p.id AND r.reaction_type = 'like'
			LEFT JOIN comments c ON c.post_id = p.id
			WHERE p.id = ?
			GROUP BY p.id, p.title, p.content, p.category_id, p.user_id, u.username, u.avatar_url, p.image_url, p.created_at
		`

		var post Post
		err := db.QueryRow(query, postID).Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CategoryID,
			&post.AuthorID,
			&post.AuthorName,
			&post.AvatarURL,
			&post.ImageURL,
			&post.CreatedAt,
			&post.LikesCount,
			&post.CommentsCount,
		)

		if err != nil {
			log.Println("❌ Post non trouvé:", err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Post non trouvé",
			})
			return
		}

		// Récupérer les tags
		tagRows, err := db.Query(`
			SELECT t.name 
			FROM tags t
			JOIN post_tags pt ON t.id = pt.tag_id
			WHERE pt.post_id = ?
		`, post.ID)

		if err == nil {
			tags := []string{}
			for tagRows.Next() {
				var tag string
				tagRows.Scan(&tag)
				tags = append(tags, tag)
			}
			post.Tags = tags
			tagRows.Close()
		}

		log.Println("✅ Post trouvé:", post.Title)

		// Renvoyer le post
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"post":   post,
		})
	}
}
