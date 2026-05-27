// POST DETAIL PAGE
const urlParams = new URLSearchParams(window.location.search);
const postId = urlParams.get('id');

if (!postId) {
    alert('❌ Post non trouvé');
    window.location.href = '/';
}

let currentUser = null;
let commentImage = null;

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function refreshUser() {
    const userStr = localStorage.getItem('user');
    currentUser = userStr ? JSON.parse(userStr) : null;
    return currentUser;
}

async function loadPost() {
    try {
        const response = await fetch(`/api/posts/get?id=${postId}`);
        const data = await response.json();

        if (data.status === 'success' && data.post) {
            displayPost(data.post);
        }
    } catch (error) {
        console.error('❌ Erreur:', error);
    }
}

function displayPost(post) {
    const categoryNames = {
        '1': 'Manga', '2': 'Anime', '3': 'Autres œuvres',
        '4': 'L\'artiste Araki', '5': 'Communauté'
    };

    const date = new Date(post.created_at);
    const formattedDate = date.toLocaleDateString('fr-FR', {
        day: 'numeric', month: 'long', year: 'numeric',
        hour: '2-digit', minute: '2-digit'
    });

    const imageHTML = post.image_url ? 
        `<img src="${post.image_url}" alt="${escapeHtml(post.title)}" class="post-image-detail">` : '';

    const tagsHTML = post.tags && post.tags.length > 0 ?
        `<div class="post-tags-detail">${post.tags.map(tag => `<span class="post-tag-detail">${escapeHtml(tag)}</span>`).join('')}</div>` : '';

    document.getElementById('postDetail').innerHTML = `
        <div class="post-header">
            <h1 class="post-title-detail">${escapeHtml(post.title)}</h1>
            <div class="post-metadata">
                <img src="${post.avatar_url}" alt="${escapeHtml(post.author_name)}" class="author-avatar">
                <span class="post-author-detail">@${escapeHtml(post.author_name)}</span>
                <span>•</span>
                <span class="post-date-detail">${formattedDate}</span>
                <span class="post-category-badge">${categoryNames[post.category_id] || 'Autre'}</span>
            </div>
            ${tagsHTML}
        </div>
        ${imageHTML}
        <div class="post-content-detail">${escapeHtml(post.content)}</div>
        <div class="post-actions">
            <button class="action-btn" id="likeBtn">
                <span class="action-icon">❤️</span>
                <span id="likeCount">${post.likes_count || 0}</span>
            </button>
            <button class="action-btn">
                <span class="action-icon">💬</span>
                <span>${post.comments_count || 0}</span>
            </button>
        </div>
    `;

    loadComments();
}

async function loadComments() {
    refreshUser();
    
    try {
        let url = `/api/comments?post_id=${postId}`;
        if (currentUser) {
            url += `&user_email=${encodeURIComponent(currentUser.email)}`;
        }
        
        const response = await fetch(url);
        const data = await response.json();

        if (data.status === 'success') {
            displayComments(data.comments || []);
        }
    } catch (error) {
        console.error('❌ Erreur:', error);
    }
}

function displayComments(comments) {
    const commentsList = document.getElementById('commentsList');

    if (!comments || comments.length === 0) {
        commentsList.innerHTML = '<div class="no-comments">Aucun commentaire pour le moment. Soyez le premier à commenter !</div>';
        return;
    }

    commentsList.innerHTML = comments.map(comment => {
        const date = new Date(comment.created_at);
        const formattedDate = date.toLocaleDateString('fr-FR', {
            day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit'
        });

        const imageHtml = comment.image_url ? `<img src="${comment.image_url}" alt="Image" class="comment-image">` : '';

        return `
            <div class="comment-item" data-comment-id="${comment.id}">
                <div class="comment-header">
                    <div class="comment-author-info">
                        <img src="${comment.avatar_url}" alt="${escapeHtml(comment.author_name)}" class="comment-avatar">
                        <span class="comment-author">@${escapeHtml(comment.author_name)}</span>
                    </div>
                    <span class="comment-date">${formattedDate}</span>
                </div>
                <div class="comment-content">${escapeHtml(comment.content)}</div>
                ${imageHtml}
                <div class="comment-actions">
                    <button class="comment-action-btn like-btn" data-comment-id="${comment.id}" ${!currentUser ? 'disabled' : ''}>
                        <span class="like-icon">${comment.likes_count > 0 ? '❤️' : '🤍'}</span>
                        <span class="like-count">${comment.likes_count || 0}</span>
                    </button>
                    <button class="comment-action-btn favorite-btn" data-comment-id="${comment.id}" ${!currentUser ? 'disabled' : ''}>
                        <span class="favorite-icon">${comment.is_favorite ? '⭐' : '☆'}</span>
                    </button>
                </div>
            </div>
        `;
    }).join('');

    attachCommentActions();
}

function attachCommentActions() {
    document.querySelectorAll('.like-btn').forEach(btn => {
        btn.addEventListener('click', likeComment);
    });
    document.querySelectorAll('.favorite-btn').forEach(btn => {
        btn.addEventListener('click', toggleFavorite);
    });
}

async function likeComment(e) {
    e.preventDefault();
    refreshUser();
    
    if (!currentUser) {
        alert('❌ Connectez-vous pour liker');
        return;
    }

    const commentId = this.getAttribute('data-comment-id');
    
    try {
        const response = await fetch('/api/comments/like', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                comment_id: commentId,
                author_email: currentUser.email
            })
        });

        if (response.ok) {
            loadComments();
        }
    } catch (error) {
        console.error('❌ Erreur:', error);
    }
}

async function toggleFavorite(e) {
    e.preventDefault();
    refreshUser();
    
    if (!currentUser) {
        alert('❌ Connectez-vous pour ajouter en favori');
        return;
    }

    const commentId = this.getAttribute('data-comment-id');
    
    try {
        const response = await fetch('/api/favorites', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                comment_id: commentId,
                author_email: currentUser.email
            })
        });

        if (response.ok) {
            loadComments();
        }
    } catch (error) {
        console.error('❌ Erreur:', error);
    }
}

// Gestion image
document.addEventListener('DOMContentLoaded', () => {
    refreshUser();
    
    // Afficher/masquer formulaire
    if (currentUser) {
        document.getElementById('commentForm').style.display = 'block';
        document.getElementById('loginPrompt').style.display = 'none';
    } else {
        document.getElementById('commentForm').style.display = 'none';
        document.getElementById('loginPrompt').style.display = 'block';
    }

    // Upload image
    const imageInput = document.getElementById('commentImageInput');
    if (imageInput) {
        imageInput.addEventListener('change', (e) => {
            const file = e.target.files[0];
            if (!file) return;

            if (file.size > 2 * 1024 * 1024) {
                alert('❌ Max 2MB');
                imageInput.value = '';
                return;
            }

            const reader = new FileReader();
            reader.onload = (event) => {
                const img = new Image();
                img.onload = () => {
                    const canvas = document.createElement('canvas');
                    let width = img.width;
                    let height = img.height;
                    
                    if (width > 800) {
                        const ratio = 800 / width;
                        width = 800;
                        height = height * ratio;
                    }
                    
                    canvas.width = width;
                    canvas.height = height;
                    canvas.getContext('2d').drawImage(img, 0, 0, width, height);
                    
                    commentImage = canvas.toDataURL('image/jpeg', 0.8);
                    document.getElementById('previewImg').src = commentImage;
                    document.getElementById('imagePreview').style.display = 'flex';
                };
                img.src = event.target.result;
            };
            reader.readAsDataURL(file);
        });
    }

    // Supprimer image
    const removeBtn = document.getElementById('removeImageBtn');
    if (removeBtn) {
        removeBtn.addEventListener('click', (e) => {
            e.preventDefault();
            commentImage = null;
            imageInput.value = '';
            document.getElementById('imagePreview').style.display = 'none';
        });
    }

    // Soumettre commentaire
    const commentForm = document.getElementById('commentForm');
    if (commentForm) {
        commentForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            refreshUser();
            
            if (!currentUser) {
                alert('❌ Connectez-vous');
                return;
            }

            const content = document.getElementById('commentContent').value.trim();
            if (!content) {
                alert('❌ Commentaire vide');
                return;
            }

            try {
                const response = await fetch('/api/comments', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        post_id: postId,
                        author_email: currentUser.email,
                        content: content,
                        image_url: commentImage || ''
                    })
                });

                const data = await response.json();
                if (data.status === 'success') {
                    document.getElementById('commentContent').value = '';
                    commentImage = null;
                    imageInput.value = '';
                    document.getElementById('imagePreview').style.display = 'none';
                    loadComments();
                }
            } catch (error) {
                console.error('❌ Erreur:', error);
            }
        });
    }

    // Charger posts
    loadPost();
    loadRelatedPosts();
});

async function loadRelatedPosts() {
    try {
        const response = await fetch('/api/posts/list?limit=10');
        const data = await response.json();
        if (data.status === 'success' && data.posts) {
            displayRelatedPosts(data.posts);
        }
    } catch (error) {
        console.error('❌ Erreur:', error);
    }
}

function displayRelatedPosts(posts) {
    const relatedPosts = document.getElementById('relatedPosts');
    if (!posts || posts.length === 0) {
        relatedPosts.innerHTML = '<div class="loading">Aucun post</div>';
        return;
    }

    const categoryNames = {
        '1': 'Manga', '2': 'Anime', '3': 'Autres',
        '4': 'Araki', '5': 'Communauté'
    };

    relatedPosts.innerHTML = posts.slice(0, 10).map(post => `
        <a href="/post.html?id=${post.id}" class="related-post-item">
            <div class="related-post-title">${escapeHtml(post.title)}</div>
            <div class="related-post-author">@${escapeHtml(post.author_name)}</div>
            <div class="related-post-category">${categoryNames[post.category_id] || 'Autre'}</div>
        </a>
    `).join('');
}