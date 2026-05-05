// Switch between Login and Register tabs
const authTabs = document.querySelectorAll('.auth-tab');
const authForms = document.querySelectorAll('.auth-form');

authTabs.forEach(tab => {
    tab.addEventListener('click', () => {
        const targetTab = tab.getAttribute('data-tab');
        
        authTabs.forEach(t => t.classList.remove('active'));
        authForms.forEach(f => f.classList.remove('active'));
        
        tab.classList.add('active');
        const targetForm = document.getElementById(targetTab + 'Form');
        if (targetForm) {
            targetForm.classList.add('active');
        }
    });
});

// Avatar Preview
const avatarInput = document.getElementById('avatar');
const avatarPreview = document.getElementById('avatarPreview');

if (avatarInput) {
    avatarInput.addEventListener('change', (e) => {
        const file = e.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (event) => {
                avatarPreview.innerHTML = `<img src="${event.target.result}" alt="Avatar">`;
            };
            reader.readAsDataURL(file);
        }
    });
}

// Password Strength Checker
const passwordInput = document.getElementById('register-password');
const strengthBar = document.querySelector('.strength-bar');

if (passwordInput && strengthBar) {
    passwordInput.addEventListener('input', () => {
        const password = passwordInput.value;
        const strength = calculatePasswordStrength(password);
        
        strengthBar.classList.remove('weak', 'medium', 'strong');
        
        if (password.length === 0) {
            strengthBar.style.width = '0%';
        } else if (strength < 40) {
            strengthBar.classList.add('weak');
        } else if (strength < 70) {
            strengthBar.classList.add('medium');
        } else {
            strengthBar.classList.add('strong');
        }
    });
}

function calculatePasswordStrength(password) {
    let strength = 0;
    if (password.length >= 8) strength += 25;
    if (password.length >= 12) strength += 15;
    if (/[a-z]/.test(password)) strength += 15;
    if (/[A-Z]/.test(password)) strength += 15;
    if (/\d/.test(password)) strength += 15;
    if (/[^a-zA-Z0-9]/.test(password)) strength += 15;
    return strength;
}

const loginForm = document.getElementById('loginForm');

if (loginForm) {
    loginForm.addEventListener('submit', (e) => {
        e.preventDefault();
        
        const formData = new FormData(loginForm);
        const email = formData.get('email');
        const password = formData.get('password');
        
        console.log('Login attempt:', email);
        
        // Créer utilisateur temporaire
        const user = {
            username: email.split('@')[0],
            email: email,
            avatar: 'images/default-avatar.png',
            bio: '',
            location: '',
            favorite_jojo_part: '',
            favorite_stand: ''
        };
        localStorage.setItem('user', JSON.stringify(user));
        
        console.log('User saved:', user);
        
        window.location.href = 'profile.html';
    });
}

const registerForm = document.getElementById('registerForm');

if (registerForm) {
    registerForm.addEventListener('submit', (e) => {
        e.preventDefault();
        
        console.log('Register form submitted');
        
        const formData = new FormData(registerForm);
        
        const password = formData.get('password');
        const confirm = formData.get('confirm');
        
        if (password !== confirm) {
            alert('Les mots de passe ne correspondent pas !');
            return;
        }
        
        if (password.length < 8) {
            alert('Le mot de passe doit contenir au moins 8 caractères !');
            return;
        }
        
        const username = formData.get('username');
        const email = formData.get('email');
        
        if (!username || !email) {
            alert('Veuillez remplir tous les champs obligatoires !');
            return;
        }
        
        console.log('Validation passed');
        
        const interests = [];
        formData.getAll('interests').forEach(interest => {
            interests.push(interest);
        });
        
        console.log('Interests:', interests);
        
        const newUser = {
            username: username,
            email: email,
            password: password,
            favorite_jojo_part: formData.get('favorite_jojo_part') || '',
            favorite_stand: formData.get('favorite_stand') || '',
            interests: interests,  // ✅ Maintenant c'est défini !
            bio: '',
            location: '',
            avatar: 'images/default-avatar.png'
        };
        
        console.log('User object created:', newUser);
        
        const avatarFile = formData.get('avatar');
        
        if (avatarFile && avatarFile.size > 0) {
            console.log('Avatar file detected, reading...');
            const reader = new FileReader();
            reader.onload = (event) => {
                newUser.avatar = event.target.result;
                saveUserAndRedirect(newUser);
            };
            reader.readAsDataURL(avatarFile);
        } else {
            console.log('No avatar, using default');
            saveUserAndRedirect(newUser);
        }
    });
}

function saveUserAndRedirect(user) {
    console.log('Envoi au backend:', user);
    fetch('/api/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'success') {
            localStorage.setItem('user', JSON.stringify(user));
            window.location.href = 'profile.html';
        } else {
            alert('❌ Erreur : ' + data.error);
        }
    })
    .catch(error => {
        console.error('Erreur:', error);
        alert('❌ Erreur de connexion au serveur');
    });
}
console.log('%c⭐ AUTH.JS CHARGÉ ⭐', 'color: #FF6B35; font-size: 16px; font-weight: bold;');