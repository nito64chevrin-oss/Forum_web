function updateThemeIcon(theme) {
    const icon = themeToggle.querySelector('.theme-icon');
    icon.textContent = theme === 'dark' ? '☀️' : '🌙';
}

document.addEventListener('click', (e) => {
    if (!e.target.closest('.nav-item')) {
        document.querySelectorAll('.nav-item.open').forEach(i => i.classList.remove('open'));
    }
});

document.querySelectorAll('.nav-link').forEach(link => {
    link.addEventListener('click', (e) => {
        const parent   = link.closest('.nav-item');
        const dropdown = parent.querySelector('.dropdown');
        if (!dropdown) return;
        e.preventDefault();
        const isOpen = parent.classList.contains('open');
        document.querySelectorAll('.nav-item.open').forEach(i => i.classList.remove('open'));
        if (!isOpen) parent.classList.add('open');
    });
});

const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.classList.add('visible');
            observer.unobserve(entry.target);
        }
    });
}, { threshold: 0.06 });

document.querySelectorAll('.post-card, .category-card, .comment-card').forEach(el => observer.observe(el));

document.querySelectorAll('.tab-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        const group = btn.closest('.tabs');
        group.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');

        const target = btn.dataset.tab;
        document.querySelectorAll('.tab-panel').forEach(panel => {
            panel.classList.toggle('tab-panel--active', panel.dataset.panel === target);
        });
    });
});

document.querySelectorAll('.reaction-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        btn.classList.toggle('active');
        const count = btn.querySelector('.reaction-count');
        if (!count) return;
        const n = parseInt(count.textContent);
        count.textContent = btn.classList.contains('active') ? n + 1 : n - 1;
    });
});

document.querySelectorAll('.filter-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        btn.closest('.filter-group').querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
    });
});

document.querySelectorAll('.tag-checkbox').forEach(cb => {
    cb.addEventListener('change', () => {
        const label = cb.nextElementSibling;
        const color = label.dataset.color;
        if (cb.checked) {
            label.style.background = color;
            label.style.borderColor = color;
        } else {
            label.style.background = '';
            label.style.borderColor = '';
        }
    });
});
