// Espera a que el DOM esté completamente cargado antes de ejecutar el script 
document.addEventListener("DOMContentLoaded", () => {
    // Botón de Gestión de recomendaciones
    document.getElementById("recommendations").addEventListener("click", () => {
        window.location.href = "/manage-recommendations"; // Cambia esta ruta según backend
    });

    // Botón de Gestión de comentarios
    document.getElementById("comments").addEventListener("click", () => {
        window.location.href = "/manage-comments"; // Cambia esta ruta según backend
    });

    // Botón de Gestión de contenido
    document.getElementById("content").addEventListener("click", () => {
        window.location.href = "/manage-content"; // Cambia esta ruta según backend
    });

    // Botón de Gestión de usuarios
    document.getElementById("users").addEventListener("click", () => {
        window.location.href = "/manage-users"; // Cambia esta ruta según backend
    });

    // Botón de Creación de reportes
    document.getElementById("reports").addEventListener("click", () => {
        window.location.href = "/create-reports"; // Cambia esta ruta según backend
    });

    // Botón de Logout
    document.querySelector(".logout").addEventListener("click", () => {
        if (confirm("¿Estás seguro de que deseas cerrar sesión?")) {
            fetch("/api/logout", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
            })
                .then((response) => {
                    if (response.ok) {
                        alert("Sesión cerrada exitosamente.");
                        window.location.href = "/";
                    } else {
                        alert("Error cerrando sesión.");
                    }
                })
                .catch((err) => console.error("Error:", err));
        }
    });
});
