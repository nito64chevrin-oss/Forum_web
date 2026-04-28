// Theme Toggle with JoJo-style animation
const themeToggle = document.getElementById('themeToggle');
const html = document.documentElement;

// Check for saved theme preference or default to 'dark'
const currentTheme = localStorage.getItem('theme') || 'dark';
html.setAttribute('data-theme', currentTheme);
updateThemeIcon(currentTheme);

themeToggle.addEventListener('click', () => {
    const currentTheme = html.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    
    // JoJo-style color flash animation
    jojoThemeTransition(newTheme);
});

function jojoThemeTransition(newTheme) {
    const flash = document.createElement('div');
    flash.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        pointer-events: none;
        z-index: 9999;
        background: linear-gradient(45deg, 
            #4A5A9E 0%, 
            #6A5A9E 25%, 
            #5BA8B8 50%, 
            #C85A88 75%, 
            #4A5A9E 100%);
        opacity: 0;
        transition: opacity 0.1s ease-in;
    `;
    document.body.appendChild(flash);
    
    // Trigger flash
    setTimeout(() => {
        flash.style.opacity = '0.3';
    }, 10);
    
    // Flash peak
    setTimeout(() => {
        flash.style.opacity = '0';
        flash.style.transition = 'opacity 0.3s ease-out';
    }, 100);
    
    // Change theme during flash
    setTimeout(() => {
        html.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
        updateThemeIcon(newTheme);
    }, 100);
    
    // Remove flash element
    setTimeout(() => {
        document.body.removeChild(flash);
    }, 500);
}

function updateThemeIcon(theme) {
    const icon = themeToggle.querySelector('.theme-icon');
    icon.textContent = theme === 'dark' ? '🌙' : '☀️';
}

// Tabs functionality
const tabs = document.querySelectorAll('.tab');
tabs.forEach(tab => {
    tab.addEventListener('click', () => {
        // Remove active class from all tabs
        tabs.forEach(t => t.classList.remove('active'));
        // Add active class to clicked tab
        tab.classList.add('active');
        
        // Here you would load different content based on tab
        // For now, just visual feedback
        console.log('Tab clicked:', tab.textContent);
    });
});

// Categories toggle (sidebar button)
const categoriesToggle = document.getElementById('categoriesToggle');
if (categoriesToggle) {
    categoriesToggle.addEventListener('click', () => {
        const arrow = categoriesToggle.querySelector('.arrow');
        
        // Rotate arrow animation
        if (arrow.style.transform === 'rotate(90deg)') {
            arrow.style.transform = 'rotate(0deg)';
        } else {
            arrow.style.transform = 'rotate(90deg)';
        }
        
        // Here you would show/hide categories sidebar
        console.log('Categories toggle clicked');
    });
}

// Category items click handling
const categoryItems = document.querySelectorAll('.category-item');
categoryItems.forEach(item => {
    item.addEventListener('click', (e) => {
        // Don't navigate if clicking on a topic link
        if (e.target.closest('.topic-link') || e.target.closest('.tag')) {
            return;
        }
        
        const category = item.getAttribute('data-category');
        console.log('Category clicked:', category);
        // Here you would navigate to category page
    });
});

// Topic links
const topicLinks = document.querySelectorAll('.topic-link');
topicLinks.forEach(link => {
    link.addEventListener('click', (e) => {
        e.preventDefault();
        e.stopPropagation();
        console.log('Topic clicked:', link.querySelector('.topic-title').textContent);
        // Here you would navigate to topic page
    });
});

// Tags click handling
const tags = document.querySelectorAll('.tag');
tags.forEach(tag => {
    tag.addEventListener('click', (e) => {
        e.stopPropagation();
        console.log('Tag clicked:', tag.textContent);
        // Here you would filter by tag
    });
});

// Smooth scroll behavior
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// Add loading animation on page load
window.addEventListener('load', () => {
    document.body.style.opacity = '0';
    document.body.style.transition = 'opacity 0.3s ease';
    
    setTimeout(() => {
        document.body.style.opacity = '1';
    }, 100);
});

// Keyboard shortcuts
document.addEventListener('keydown', (e) => {
    // Ctrl/Cmd + K for search
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        document.querySelector('.btn-search')?.click();
        console.log('Search shortcut activated');
    }
    
    // Ctrl/Cmd + D for dark mode toggle
    if ((e.ctrlKey || e.metaKey) && e.key === 'd') {
        e.preventDefault();
        themeToggle.click();
    }
});

// Mobile menu toggle
const menuBtn = document.querySelector('.btn-menu');
if (menuBtn) {
    menuBtn.addEventListener('click', () => {
        console.log('Mobile menu toggled');
        // Here you would show/hide mobile menu
    });
}

// Intersection Observer for scroll animations
const observerOptions = {
    threshold: 0.1,
    rootMargin: '0px 0px -50px 0px'
};

const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.style.opacity = '0';
            entry.target.style.transform = 'translateY(20px)';
            entry.target.style.transition = 'opacity 0.5s ease, transform 0.5s ease';
            
            setTimeout(() => {
                entry.target.style.opacity = '1';
                entry.target.style.transform = 'translateY(0)';
            }, 100);
            
            observer.unobserve(entry.target);
        }
    });
}, observerOptions);

// Observe category items for scroll animation
document.querySelectorAll('.category-item').forEach(item => {
    observer.observe(item);
});

// Console easter egg
console.log('%c⭐ ARAKI FORUM ⭐', 'color: #FF6B35; font-size: 20px; font-weight: bold;');
console.log('%cBienvenue sur le forum dédié à Hirohiko Araki !', 'color: #5BA8B8; font-size: 14px;');
console.log('%cRaccourcis clavier:', 'color: #4A5A9E; font-weight: bold;');
console.log('%c  • Ctrl/Cmd + K : Recherche', 'color: #999;');
console.log('%c  • Ctrl/Cmd + D : Toggle Dark/Light mode', 'color: #999;');
