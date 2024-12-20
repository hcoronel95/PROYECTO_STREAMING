document.addEventListener("DOMContentLoaded", () => {
    // Gestión de contenido 
    document.getElementById("content").addEventListener("click", () => {
        window.location.href = "/pages/content-management.html";
    });

    // Gestión de usuarios
    document.getElementById("users").addEventListener("click", () => {
        window.location.href = "/pages/desarrollo.html";
    });

    // Creación de reportes
    document.getElementById("reports").addEventListener("click", () => {
        window.location.href = "/pages/desarrollo.html";
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