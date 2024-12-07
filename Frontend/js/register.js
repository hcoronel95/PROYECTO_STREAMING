document.addEventListener('DOMContentLoaded', function() {
    const registerForm = document.getElementById('registerForm');
    const errorMessage = document.getElementById('errorMessage');

    registerForm?.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        try {
            const response = await fetch('/api/users/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ name, email, password })
            });

            if (response.ok) {
                const data = await response.json();
                alert('Registro exitoso. Por favor inicia sesión.');
                window.location.href = '/pages/login.html';
            } else {
                errorMessage.style.display = 'block';
                errorMessage.textContent = 'Error en el registro. Por favor intenta nuevamente.';
            }
        } catch (error) {
            console.error('Error:', error);
            errorMessage.style.display = 'block';
            errorMessage.textContent = 'Error al conectar con el servidor';
        }
    });
});