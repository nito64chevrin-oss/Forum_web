// GET les params de l'URL
const urlParams = new URLSearchParams(window.location.search);
const categoryId = urlParams.get('category');
const selectedTag = urlParams.get('tag') || null;

const categoryNames = {
    '1': 'Manga',
    '2': 'Anime',
    '3': 'Autres œuvres',
    '4': 'Hirohiko Araki',
    '5': 'Communauté'
};

let allPosts = [];
let allTags = [];
let currentFilter = selectedTag;

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Charger les posts de la catégorie
async function loadCategoryPosts() {
    try {
        const response = await fetch(`/api/posts/list?category_id=${categoryId}`);
        const data = await response.json();

        if (data.status === 'success') {
            allPosts = data.posts || [];
            extractTags();
            displayTags();
            filterAndDisplay();
        }
    } catch (error) {
        console.error('❌ Erreur:', error);
    }
}

// Extraire tous les tags uniques
function extractTags() {
    const tagsSet = new Set();
    allPosts.forEach(post => {
        if (post.tags && Array.isArray(post.tags)) {
            post.tags.forEach(tag => tagsSet.add(tag));
        }
    });
    allTags = Array.from(tagsSet).sort();
}

// Afficher les tags filtrables
function displayTags() {
    const tagsList = document.getElementById('tagsList');
    
    const allTagsHtml = `
        <button class="tag-filter-btn ${!currentFilter ? 'active' : ''}" data-tag="all">
            Tous
        </button>
    `;

    const tagsHtml = allTags.map(tag => `
        <button class="tag-filter-btn ${currentFilter === tag ? 'active' : ''}" data-tag="${escapeHtml(tag)}">
            ${escapeHtml(tag)}
        </button>
    `).join('');

    tagsList.innerHTML = allTagsHtml + tagsHtml;

    // Ajouter événements
    document.querySelectorAll('.tag-filter-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            currentFilter = btn.dataset.tag === 'all' ? null : btn.dataset.tag;
            
            // Mettre à jour l'URL
            const params = new URLSearchParams();
            params.set('category', categoryId);
            if (currentFilter) params.set('tag', currentFilter);
            window.history.pushState({}, '', `?${params.toString()}`);
            
            filterAndDisplay();
            displayTags();
        });
    });
}

// Filtrer et afficher les posts
function filterAndDisplay() {
    let filtered = allPosts;

    if (currentFilter) {
        filtered = allPosts.filter(post => 
            post.tags && post.tags.includes(currentFilter)
        );
    }

    displayPosts(filtered);
}

// Afficher les posts
function displayPosts(posts) {
    const postsGrid = document.getElementById('postsGrid');

    if (!posts || posts.length === 0) {
        postsGrid.innerHTML = '<div style="grid-column: 1/-1; text-align: center; padding: 2rem; color: var(--text-muted);">Aucun post trouvé</div>';
        return;
    }

    postsGrid.innerHTML = posts.map(post => {
        const date = new Date(post.created_at);
        const formattedDate = date.toLocaleDateString('fr-FR', {
            day: 'numeric',
            month: 'short'
        });

        const tagsHtml = post.tags && post.tags.length > 0 ?
            `<div class="post-card-tags">${post.tags.slice(0, 3).map(tag => 
                `<a href="?category=${categoryId}&tag=${encodeURIComponent(tag)}" class="post-card-tag">${escapeHtml(tag)}</a>`
            ).join('')}</div>` : '';

        return `
            <a href="/post.html?id=${post.id}" class="post-card">
                <div class="post-card-header">
                    <h3 class="post-card-title">${escapeHtml(post.title)}</h3>
                </div>
                <div class="post-card-meta">
                    <span class="post-card-author">@${escapeHtml(post.author_name)}</span>
                    <span class="post-card-date">${formattedDate}</span>
                </div>
                ${tagsHtml}
                <div class="post-card-stats">
                    <span>❤️ ${post.likes_count || 0}</span>
                    <span>💬 ${post.comments_count || 0}</span>
                </div>
            </a>
        `;
    }).join('');
}

// Au chargement
document.addEventListener('DOMContentLoaded', () => {
    // Titre catégorie
    document.getElementById('categoryTitle').textContent = categoryNames[categoryId] || 'Catégorie';
    
    // Charger les posts
    loadCategoryPosts();
});

console.log('📂 Catégorie:', categoryId, '- Tag filter:', selectedTag);