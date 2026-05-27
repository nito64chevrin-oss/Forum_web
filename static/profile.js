

let currentUser = null;

// Charger l'utilisateur depuis localStorage
function loadCurrentUser() {
    const userStr = localStorage.getItem('user');
    if (!userStr) {
        return null;
    }
    try {
        currentUser = JSON.parse(userStr);
        return currentUser;
    } catch (e) {
        console.error('Error parsing user:', e);
        return null;
    }
}

document.addEventListener('DOMContentLoaded', () => {
    currentUser = loadCurrentUser();

    if (!currentUser) {
        window.location.href = 'auth.html';
        return;
    }

    loadUserProfile();
});

// Charger le profil depuis l'API
async function loadUserProfile() {
    try {
        const response = await fetch(`/api/user?email=${encodeURIComponent(currentUser.email)}`);
        const data = await response.json();

        if (data.status === 'success') {
            displayUserProfile(data.user);
            if (data.stats) {
                displayUserStats(data.stats);
            }
        } else {
            console.error('API Error:', data.error);
        }
    } catch (error) {
        console.error('Fetch Error:', error);
    }
}

// Afficher le profil
function displayUserProfile(user) {
    
    // Avatar - correcter: c'est #currentAvatar pas #profileAvatar
    const avatarImg = document.getElementById('currentAvatar');
    if (avatarImg) {
        avatarImg.src = user.avatar_url || 'static/images/default-avatar.png';
    }

    // Infos principales
    const usernameField = document.getElementById('profileUsername');
    const emailField = document.getElementById('profileEmail');
    
    if (usernameField) usernameField.textContent = user.username;
    if (emailField) emailField.textContent = user.email;

    // Date de création - dans le HTML c'est .member-since
    const memberSinceField = document.querySelector('.member-since');
    if (memberSinceField && user.created_at) {
        const date = new Date(user.created_at).toLocaleDateString('fr-FR', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
        memberSinceField.textContent = `Membre depuis ${date}`;
    }

    // Update form fields for editing - avec les bons IDs du HTML
    const bioInput = document.getElementById('bio');
    const standInput = document.getElementById('favorite-stand'); 
    const partInput = document.getElementById('favorite-part');    
    const locationInput = document.getElementById('location');
    
    if (bioInput) bioInput.value = user.bio || '';
    if (standInput) standInput.value = user.favorite_stand || '';
    if (partInput) partInput.value = user.favorite_jojo_part || '';
    if (locationInput) locationInput.value = user.location || '';
    
}

// Afficher les stats
function displayUserStats(stats) {
    
    const postsCount = document.getElementById('postsCount');
    const commentsCount = document.getElementById('commentsCount');
    const likesReceived = document.getElementById('likesReceived');  // ✅ correct ID

    if (postsCount) postsCount.textContent = stats.posts_count || 0;
    if (commentsCount) commentsCount.textContent = stats.comments_count || 0;
    if (likesReceived) likesReceived.textContent = stats.likes_received || 0;
    
}

// Soumettre le formulaire de profil
const profileForm = document.getElementById('profileForm');
if (profileForm) {
    profileForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const bio = document.getElementById('bio')?.value.trim() || '';
        const favoriteStand = document.getElementById('favorite-stand')?.value.trim() || '';
        const favoritePart = document.getElementById('favorite-part')?.value.trim() || '';
        const location = document.getElementById('location')?.value.trim() || '';

        try {
            const response = await fetch('/api/user/update', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    email: currentUser.email,
                    username: currentUser.username,
                    bio: bio,
                    favorite_stand: favoriteStand,
                    favorite_jojo_part: favoritePart,
                    location: location
                })
            });

            const data = await response.json();
            if (data.status === 'success') {
                alert('✅ Profil mis à jour !');
                loadUserProfile();
            } else {
                alert('❌ Erreur: ' + data.error);
            }
        } catch (error) {
            console.error('Error:', error);
            alert('❌ Erreur: ' + error.message);
        }
    });
}


const avatarInput = document.getElementById('newAvatar');
if (avatarInput) {
    avatarInput.addEventListener('change', async (e) => {
        const file = e.target.files[0];
        if (!file) return;


        // Vérifier la taille (max 2MB)
        if (file.size > 2 * 1024 * 1024) {
            alert('❌ Image trop volumineuse ! Max 2MB');
            avatarInput.value = '';
            return;
        }

        // Lire et compresser l'image
        const reader = new FileReader();
        reader.onload = (event) => {
            const img = new Image();
            img.onload = () => {
                const canvas = document.createElement('canvas');
                canvas.width = 200;
                canvas.height = 200;
                
                const ctx = canvas.getContext('2d');
                ctx.drawImage(img, 0, 0, 200, 200);
                
                const compressedImage = canvas.toDataURL('image/jpeg', 0.8);
                
                // Afficher l'aperçu
                const avatarImg = document.getElementById('currentAvatar');
                if (avatarImg) {
                    avatarImg.src = compressedImage;
                }

                // Sauvegarder immédiatement
                saveAvatar(compressedImage);
            };
            img.src = event.target.result;
        };
        reader.readAsDataURL(file);
    });
}

async function saveAvatar(imageData) {
    try {
        const response = await fetch('/api/user/update', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                email: currentUser.email,
                avatar_url: imageData
            })
        });

        const data = await response.json();
        if (data.status === 'success') {
            currentUser = data.user;
            localStorage.setItem('user', JSON.stringify(currentUser));
        } else {
        }
    } catch (error) {
    }
}