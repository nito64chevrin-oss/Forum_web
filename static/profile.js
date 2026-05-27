function loadUserProfile() {
    const tempData = localStorage.getItem('user');
    
    if (!tempData) {
        window.location.href = 'auth.html';
        return;
    }
    
    const tempUser = JSON.parse(tempData);
    const userEmail = tempUser.email;
    
    
    fetch(`/api/user?email=${userEmail}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Utilisateur non trouvé');
            }
            return response.json();
        })
        .then(data => {
            if (data.status === 'success') {
                const user = data.user;
                
                
                document.getElementById('profileUsername').textContent = user.username;
                document.getElementById('profileEmail').textContent = user.email;
                
                const currentAvatar = document.getElementById('currentAvatar');
                if (user.avatar_url && user.avatar_url.startsWith('data:image')) {
                    currentAvatar.src = user.avatar_url;
                    currentAvatar.style.display = 'block';
                } else {
                    currentAvatar.src = 'images/default-avatar.png';
                    currentAvatar.style.display = 'block';
                }
                document.getElementById('bio').value = user.bio || '';
                document.getElementById('location').value = user.location || '';
                document.getElementById('favorite-part').value = user.favorite_jojo_part || '';
                document.getElementById('favorite-stand').value = user.favorite_stand || '';

                const memberSinceElement = document.querySelector('.member-since');
                if (memberSinceElement && user.created_at) {
                    console.log('Created at:', user.created_at);
                    const date = new Date(user.created_at);
                    const options = { year: 'numeric', month: 'long', day: 'numeric' };
                    memberSinceElement.textContent = `Membre depuis ${date.toLocaleDateString('fr-FR', options)}`;
                }

                const interestsContainer = document.getElementById('userInterests');
                if (interestsContainer && user.interests && user.interests.length > 0) {
                    interestsContainer.innerHTML = '';
                    const tagColors = ['tag-blue', 'tag-cyan', 'tag-violet', 'tag-pink'];
                    
                    user.interests.forEach((interest, index) => {
                        const tag = document.createElement('span');
                        tag.className = `tag-item ${tagColors[index % 4]}`;
                        tag.textContent = interest.charAt(0).toUpperCase() + interest.slice(1);
                        interestsContainer.appendChild(tag);
                    });
                } else if (interestsContainer) {
                    interestsContainer.innerHTML = '<p style="color: var(--text-muted); font-size: 0.9rem;">Aucun centre d\'intérêt sélectionné</p>';
                }
            }
        })
        .catch(error => {
            console.error('Erreur lors du chargement du profil:', error);
            alert('Impossible de charger le profil.');
            window.location.href = 'auth.html';
        });
}

const newAvatarInput = document.getElementById('newAvatar');

if (newAvatarInput) {
    newAvatarInput.addEventListener('change', (e) => {
        const file = e.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (event) => {
                const preview = document.getElementById('profileAvatarPreview');
                preview.innerHTML = `<img src="${event.target.result}" alt="Avatar" id="currentAvatar" style="width: 120px; height: 120px; border-radius: 50%; object-fit: cover;">`;
            };
            reader.readAsDataURL(file);
        }
    });
}

// À ajouter dans profile.js

// Fonction pour afficher les stats utilisateur
function displayUserStats(stats) {
    if (!stats) return;
    
    // Stats principales uniquement
    document.getElementById('postsCount').textContent = stats.posts_count || 0;
    document.getElementById('commentsCount').textContent = stats.comments_count || 0;
    document.getElementById('likesReceived').textContent = stats.likes_received || 0;
    
    console.log('📊 Stats affichées:', stats);
}

// Modifier la fonction loadUserProfile existante pour inclure les stats
async function loadUserProfile() {
    const user = localStorage.getItem('user');
    
    if (!user) {
        window.location.href = '/auth.html';
        return;
    }
    
    const userData = JSON.parse(user);
    
    try {
        const response = await fetch(`/api/user?email=${encodeURIComponent(userData.email)}`);
        const data = await response.json();
        
        if (data.status === 'success') {
            // Afficher les infos de profil
            displayUserProfile(data.user);
            
            // Afficher les stats
            if (data.stats) {
                displayUserStats(data.stats);
            }
        } else {
            console.error('❌ Erreur:', data.error);
            alert('Erreur lors du chargement du profil');
        }
    } catch (error) {
        console.error('❌ Erreur réseau:', error);
        alert('Erreur de connexion au serveur');
    }
}

// Appeler au chargement de la page
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', loadUserProfile);
} else {
    loadUserProfile();
}

const profileForm = document.getElementById('profileForm');
if (profileForm) {
    profileForm.addEventListener('submit', (e) => {
        e.preventDefault();
        
        const formData = new FormData(profileForm);
        const userData = JSON.parse(localStorage.getItem('user'));
        
        function sendUpdate(avatarBase64) {
            const updateData = {
                email: userData.email,
                bio: formData.get('bio'),
                location: formData.get('location'),
                favorite_jojo_part: formData.get('favorite_jojo_part'),
                favorite_stand: formData.get('favorite_stand'),
                avatar_url: avatarBase64 || ''
            };
            
            fetch('/api/user/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(updateData)
            })
            .then(response => response.json())
            .then(data => {
                if (data.status === 'success') {
                    alert('Profil mis à jour !');
                    window.location.reload();
                } else {
                    alert('Erreur : ' + data.error);
                }
            })
            .catch(error => {
                console.error('Erreur:', error);
                alert('Erreur de connexion');
            });
        }
        const newAvatar = formData.get('avatar');
        if (newAvatar && newAvatar.size > 0) {
            const reader = new FileReader();
            reader.onload = (event) => {
                sendUpdate(event.target.result)
            };
            reader.readAsDataURL(newAvatar);
        } else {
            sendUpdate('');
        }
    });
}

loadUserProfile();
