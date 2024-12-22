document.addEventListener("DOMContentLoaded", () => {
    const userForm = document.getElementById("userForm");
    const userTableBody = document.getElementById("userTableBody");
    const usernameInput = document.getElementById("username");
    const emailInput = document.getElementById("email");
    const roleInput = document.getElementById("role");
    const userIdInput = document.getElementById("userId");
    const alertContainer = document.createElement("div"); // Contenedor para mensajes de alerta
    userForm.parentElement.insertBefore(alertContainer, userForm); // Insertar antes del formulario

    // Datos iniciales
    let users = [
        { username: 'Juan Perez', email: 'juan.perez@example.com', role: 'user' },
        { username: 'Ana Gomez', email: 'ana.gomez@example.com', role: 'user' },
        { username: 'Carlos Lopez', email: 'carlos.lopez@example.com', role: 'user' },
    ];

    let editMode = false;

    // Función para renderizar la tabla de usuarios
    const renderTable = () => {
        userTableBody.innerHTML = ""; // Limpiar el cuerpo de la tabla
        users.forEach((user, index) => {
            const row = document.createElement("tr");
            row.innerHTML = `
                <td>${index + 1}</td>
                <td>${user.username}</td>
                <td>${user.email}</td>
                <td>${user.role}</td>
                <td>
                    <button class="btn btn-sm btn-warning edit-user" data-index="${index}">Editar</button>
                    <button class="btn btn-sm btn-danger delete-user" data-index="${index}">Eliminar</button>
                </td>
            `;
            userTableBody.appendChild(row);
        });

        // Asignar eventos a botones de editar y eliminar
        document.querySelectorAll(".edit-user").forEach(button => {
            button.addEventListener("click", handleEditUser);
        });
        document.querySelectorAll(".delete-user").forEach(button => {
            button.addEventListener("click", handleDeleteUser);
        });
    };

    // Mostrar mensaje de alerta
    const showAlert = (message, type = "success") => {
        alertContainer.innerHTML = `
            <div class="alert alert-${type} alert-dismissible fade show" role="alert">
                ${message}
                <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
            </div>
        `;
    };

    // Función para manejar el envío del formulario
    userForm.addEventListener("submit", (e) => {
        e.preventDefault(); // Evitar recargar la página

        const username = usernameInput.value.trim();
        const email = emailInput.value.trim();
        const role = roleInput.value;
        const userId = userIdInput.value;

        if (editMode) {
            // Editar usuario existente
            const index = parseInt(userId, 10);
            users[index] = { username, email, role };
            editMode = false;
        } else {
            // Agregar nuevo usuario
            users.push({ username, email, role });
        }

        // Limpiar el formulario
        userForm.reset();
        userIdInput.value = "";

        // Renderizar la tabla y mostrar éxito
        renderTable();
        showAlert("Usuario guardado con éxito");
    });

    // Función para editar usuario
    const handleEditUser = (e) => {
        const index = e.target.dataset.index;
        const user = users[index];

        // Llenar el formulario con los datos del usuario
        usernameInput.value = user.username;
        emailInput.value = user.email;
        roleInput.value = user.role;
        userIdInput.value = index;

        editMode = true; // Activar modo edición
    };

    // Función para eliminar usuario
    const handleDeleteUser = (e) => {
        const index = e.target.dataset.index;

        // Confirmar antes de eliminar
        if (confirm("¿Estás seguro de eliminar este usuario?")) {
            users.splice(index, 1); // Eliminar usuario del array
            renderTable(); // Actualizar tabla
            showAlert("Usuario eliminado con éxito", "danger"); // Mostrar alerta
        }
    };

    // Renderizar la tabla inicialmente
    renderTable();
});
