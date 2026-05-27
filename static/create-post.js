// Vérifier si l'utilisateur est connecté
function checkAuth() {
    const user = localStorage.getItem('user');
    if (!user) {
        alert('❌ Vous devez être connecté pour créer un post !');
        window.location.href = '/auth.html';
        return null;
    }
    return JSON.parse(user);
}

const currentUser = checkAuth();

// Compteur de caractères
const contentTextarea = document.getElementById('content');
const charCount = document.getElementById('charCount');

if (contentTextarea && charCount) {
    contentTextarea.addEventListener('input', () => {
        charCount.textContent = contentTextarea.value.length;
    });
}

// Preview image
const imageInput = document.getElementById('image');
const imagePreview = document.getElementById('imagePreview');

if (imageInput && imagePreview) {
    imageInput.addEventListener('change', (e) => {
        const file = e.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (event) => {
                imagePreview.innerHTML = `<img src="${event.target.result}" alt="Preview">`;
            };
            reader.readAsDataURL(file);
        } else {
            imagePreview.innerHTML = '';
        }
    });
}

// Limitation tags (max 5)
const tagCheckboxes = document.querySelectorAll('input[name="tags"]');
const tagCountElement = document.getElementById('tagCount');

function updateTagCount() {
    const checkedCount = document.querySelectorAll('input[name="tags"]:checked').length;
    tagCountElement.textContent = `${checkedCount}/5 tags sélectionnés`;
    
    // Désactiver les checkboxes non cochées si on a atteint 5
    if (checkedCount >= 5) {
        tagCheckboxes.forEach(checkbox => {
            if (!checkbox.checked) {
                checkbox.disabled = true;
                checkbox.parentElement.style.opacity = '0.5';
            }
        });
    } else {
        tagCheckboxes.forEach(checkbox => {
            checkbox.disabled = false;
            checkbox.parentElement.style.opacity = '1';
        });
    }
}

tagCheckboxes.forEach(checkbox => {
    checkbox.addEventListener('change', updateTagCount);
});

// Soumission du formulaire
const createPostForm = document.getElementById('createPostForm');

if (createPostForm) {
    createPostForm.addEventListener('submit', (e) => {
        e.preventDefault();
        
        console.log('📝 Soumission du formulaire...');
        
        if (!currentUser) {
            alert('❌ Vous devez être connecté !');
            return;
        }
        
        const formData = new FormData(createPostForm);
        
        // Récupérer les données
        const title = formData.get('title').trim();
        const content = formData.get('content').trim();
        
        // Récupérer les catégories cochées
        const selectedCategories = [];
        document.querySelectorAll('input[name="categories"]:checked').forEach(cb => {
            selectedCategories.push(cb.value);
        });
        
        // Validation
        if (!title || !content) {
            alert('❌ Veuillez remplir tous les champs obligatoires !');
            return;
        }
        
        if (selectedCategories.length === 0) {
            alert('❌ Veuillez sélectionner au moins une catégorie !');
            return;
        }
        
        if (title.length < 5) {
            alert('❌ Le titre doit contenir au moins 5 caractères !');
            return;
        }
        
        if (content.length < 20) {
            alert('❌ Le contenu doit contenir au moins 20 caractères !');
            return;
        }
        
        // Récupérer les tags sélectionnés
        const selectedTags = [];
        formData.getAll('tags').forEach(tag => {
            selectedTags.push(tag);
        });
        
        console.log('✅ Données validées');
        console.log('Tags:', selectedTags);
        
        // Gérer l'image
        const imageFile = formData.get('image');
        
        if (imageFile && imageFile.size > 0) {
            console.log('📸 Lecture de l\'image...');
            const reader = new FileReader();
            reader.onload = (event) => {
                createPost(title, content, selectedCategories[0], selectedCategories, selectedTags, event.target.result);
            };
            reader.readAsDataURL(imageFile);
        } else {
            createPost(title, content, selectedCategories[0], selectedCategories, selectedTags, '');
        }
    });
}

function createPost(title, content, categoryID, categoryIDs, tags, imageBase64) {
    console.log('🚀 Envoi au backend...');
    
    const postData = {
        title: title,
        content: content,
        category_id: categoryID,        // Catégorie principale
        category_ids: categoryIDs,      // Toutes les catégories
        author_email: currentUser.email,
        image_url: imageBase64,
        tags: tags
    };
    
    console.log('📦 Données envoyées:', postData);
    console.log('📂 Catégories:', categoryIDs.join(', '));
    
    fetch('/api/posts', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(postData)
    })
    .then(response => response.json())
    .then(data => {
        console.log('Réponse reçue:', data);
        
        if (data.status === 'success') {
            console.log('✅ Post créé avec succès !');
            alert('✅ Votre post a été publié avec succès !');
            window.location.href = '/';
        } else {
            console.error('❌ Erreur:', data.error);
            alert('❌ Erreur : ' + data.error);
        }
    })
    .catch(error => {
        console.error('❌ Erreur réseau:', error);
        alert('❌ Erreur de connexion au serveur');
    });
}
