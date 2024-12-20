document.addEventListener('DOMContentLoaded', () => {
    const uploadForm = document.getElementById('uploadSongForm');
    const uploadStatus = document.getElementById('uploadStatus');

    uploadForm?.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        uploadStatus.innerHTML = '<div class="alert alert-info">Subiendo archivo...</div>';
        
        const formData = new FormData();
        const songFile = document.getElementById('songFile').files[0];
        
        if (!songFile) {
            uploadStatus.innerHTML = '<div class="alert alert-danger">Por favor selecciona un archivo</div>';
            return;
        }
        
        // Verificar tamaño del archivo (10MB = 10 * 1024 * 1024 bytes)
        if (songFile.size > 10 * 1024 * 1024) {
            uploadStatus.innerHTML = '<div class="alert alert-danger">El archivo excede el límite de 10MB</div>';
            return;
        }
        
        // Agregar todos los campos al FormData
        formData.append('songFile', songFile);
        formData.append('title', document.getElementById('title').value);
        formData.append('artist', document.getElementById('artist').value);
        formData.append('genre', document.getElementById('genre').value);
        
        try {
            const token = localStorage.getItem('userToken');
            console.log('Token usado:', token); // Para debugging
            
            const response = await fetch('/api/songs/upload', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`
                },
                body: formData
            });

            console.log('Respuesta del servidor:', response); // Para debugging
            
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Error del servidor: ${errorText}`);
            }
            
            uploadStatus.innerHTML = '<div class="alert alert-success">¡Canción subida exitosamente!</div>';
            uploadForm.reset();
        } catch (error) {
            console.error('Error completo:', error);
            uploadStatus.innerHTML = `<div class="alert alert-danger">Error: ${error.message}</div>`;
        }
    });

    // Botón de Logout
    document.querySelector('.logout')?.addEventListener('click', () => {
        if (confirm('¿Estás seguro de que deseas cerrar sesión?')) {
            fetch('/api/logout', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
            })
            .then(response => {
                if (response.ok) {
                    localStorage.removeItem('userToken');
                    window.location.href = '/';
                } else {
                    throw new Error('Error en logout');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Error al cerrar sesión');
            });
        }
    });
});