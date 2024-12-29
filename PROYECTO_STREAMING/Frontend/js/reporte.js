document.addEventListener("DOMContentLoaded", () => {
    // Obtener los elementos del DOM
    const formulario = document.getElementById("reporteForm");
    const tituloInput = document.getElementById("titulo");
    const descripcionInput = document.getElementById("descripcion");
    const statusDiv = document.getElementById("status");
    const logoutButton = document.querySelector('.logout');

    // Función para guardar el reporte
    const guardarReporte = () => {
        const titulo = tituloInput.value.trim();
        const descripcion = descripcionInput.value.trim();

        // Validar campos
        if (!titulo || !descripcion) {
            statusDiv.innerText = "Por favor, completa todos los campos.";
            statusDiv.classList.add("text-danger");
            return;
        }

        // Crear el contenido del reporte
        const contenido = `Título: ${titulo}\nDescripción: ${descripcion}\n\n`;

        // Crear un Blob con el contenido del reporte
        const blob = new Blob([contenido], { type: "text/plain" });

        // Crear un enlace temporal para descargar el archivo
        const enlace = document.createElement("a");
        enlace.href = URL.createObjectURL(blob);
        enlace.download = "reporte.txt";

        // Simular un clic para descargar el archivo
        enlace.click();

        // Actualizar el estado y limpiar el formulario
        statusDiv.innerText = "Reporte guardado con éxito.";
        statusDiv.classList.remove("text-danger");
        statusDiv.classList.add("text-success");
        formulario.reset();
    };

    // Asignar la funcionalidad de guardar reporte al botón correspondiente
    const guardarBtn = formulario.querySelector("button[type='button']");
    guardarBtn.addEventListener("click", guardarReporte);

    // Función para manejar el logout
    logoutButton?.addEventListener('click', () => {
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

