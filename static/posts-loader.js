// ========================================
// AFFICHAGE DES POSTS SUR INDEX.HTML
// ========================================

// Mapper les IDs de catégories aux slugs
const categoryMap = {
    '1': 'manga',
    '2': 'anime',
    '3': 'autres',
    '4': 'artiste',
    '5': 'communaute'
};

// Charger les posts pour toutes les catégories
function loadAllPosts() {
    Object.keys(categoryMap).forEach(categoryID => {
        loadPostsForCategory(categoryID);
    });
}

// Charger les posts pour une catégorie spécifique
function loadPostsForCategory(categoryID) {
    const categorySlug = categoryMap[categoryID];
    const container = document.querySelector(`[data-category="${categorySlug}"] .category-main`);
    
    if (!container) {
        console.log(`❌ Container non trouvé pour ${categorySlug}`);
        return;
    }
    
    // Créer la section de posts si elle n'existe pas
    let postsSection = container.querySelector('.category-posts');
    if (!postsSection) {
        postsSection = document.createElement('div');
        postsSection.className = 'category-posts';
        
        // Insérer après la description de la catégorie
        const categoryInfo = container.querySelector('.category-content');
        if (categoryInfo) {
            categoryInfo.after(postsSection);
        } else {
            container.appendChild(postsSection);
        }
    }
    
    // Afficher le skeleton de chargement
    postsSection.innerHTML = createSkeletonHTML();
    
    // Récupérer les posts
    fetch(`/api/posts/list?category_id=${categoryID}&sort=recent`)
        .then(response => response.json())
        .then(data => {
            console.log(`✅ Posts reçus pour ${categorySlug}:`, data);
            
            if (data.status === 'success' && data.posts && data.posts.length > 0) {
                displayPosts(postsSection, data.posts, categoryID);
            } else {
                postsSection.innerHTML = '<div class="no-posts">Aucun post dans cette catégorie pour le moment</div>';
            }
        })
        .catch(error => {
            console.error(`❌ Erreur chargement posts ${categorySlug}:`, error);
            postsSection.innerHTML = '<div class="no-posts">Erreur de chargement</div>';
        });
}

// Afficher les posts
function displayPosts(container, posts, categoryID) {
    // Limiter à 3 posts par catégorie sur la page d'accueil
    const displayPosts = posts.slice(0, 3);
    
    container.innerHTML = displayPosts.map(post => createPostCardHTML(post)).join('');
    
    // Ajouter un bouton "Voir plus" s'il y a plus de 3 posts
    if (posts.length > 3) {
        const viewMoreBtn = document.createElement('a');
        viewMoreBtn.href = `/category.html?id=${categoryID}`;
        viewMoreBtn.className = 'view-more-posts';
        viewMoreBtn.textContent = `Voir tous les ${posts.length} posts →`;
        container.appendChild(viewMoreBtn);
    }
}

// Créer le HTML d'une carte de post
function createPostCardHTML(post) {
    const date = new Date(post.created_at);
    const formattedDate = date.toLocaleDateString('fr-FR', { 
        day: 'numeric', 
        month: 'short'
    });
    
    const tagsHTML = post.tags && post.tags.length > 0 ? 
        `${post.tags.slice(0, 3).map(tag => `<span class="post-tag">${tag}</span>`).join('')}` : 
        '<span class="no-tags" style="font-size: 0.75rem; color: var(--text-muted); font-style: italic;">Sans tags</span>';
    
    return `
        <a href="/post.html?id=${post.id}" class="post-card">
            <h3 class="post-title">${escapeHtml(post.title)}</h3>
            <div class="post-meta">
                <span class="post-author">@${escapeHtml(post.author_name)}</span>
                <span>•</span>
                <span class="post-date">${formattedDate}</span>
            </div>
            <div class="post-tags">${tagsHTML}</div>
            <div class="post-stats">
                <div class="stat-item">
                    <span class="stat-icon">❤️</span>
                    <span>${post.likes_count || 0}</span>
                </div>
                <div class="stat-item">
                    <span class="stat-icon">💬</span>
                    <span>${post.comments_count || 0}</span>
                </div>
            </div>
        </a>
    `;
}

// Créer le skeleton de chargement
function createSkeletonHTML() {
    return `
        <div class="post-skeleton">
            <div class="skeleton-line title"></div>
            <div class="skeleton-line short"></div>
            <div class="skeleton-line text"></div>
            <div class="skeleton-line text"></div>
            <div class="skeleton-line short"></div>
        </div>
    `.repeat(3);
}

// Échapper le HTML pour éviter les injections XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Afficher le bouton "Créer un post" si connecté
function showCreatePostButton() {
    const user = localStorage.getItem('user');
    const createPostBtn = document.getElementById('createPostBtn');
    
    if (createPostBtn) {
        if (user) {
            createPostBtn.style.display = 'inline-flex';
        } else {
            createPostBtn.style.display = 'none';
        }
    }
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        loadAllPosts();
        showCreatePostButton();
    });
} else {
    loadAllPosts();
    showCreatePostButton();
}