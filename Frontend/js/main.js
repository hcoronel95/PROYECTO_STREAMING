document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    const errorMessage = document.getElementById('errorMessage');

    loginForm?.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password })
            });

            if (response.ok) {
                const data = await response.json();
                localStorage.setItem('userToken', data.token);
                
                // Redirigir según el rol del usuario
                if (data.role === 'admin') {
                    window.location.href = '/pages/admin.html';
                
                } else {
                    window.location.href = '/pages/player.html';
                }
            } else {
                errorMessage.style.display = 'block';
                errorMessage.textContent = 'Credenciales inválidas';
            }
        } catch (error) {
            console.error('Error:', error);
            errorMessage.style.display = 'block';
            errorMessage.textContent = 'Error al conectar con el servidor';
        }
    });
});