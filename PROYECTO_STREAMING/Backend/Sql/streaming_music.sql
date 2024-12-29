-- Crear la base de datos
CREATE DATABASE streaming_music;
USE streaming_music;

-- Tabla de usuarios
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'user') NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de canciones
CREATE TABLE songs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    genre VARCHAR(100) NOT NULL,
    file_size INT NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de bibliotecas de usuario
CREATE TABLE libraries (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Tabla de relación entre bibliotecas y canciones
CREATE TABLE library_songs (
    library_id INT NOT NULL,
    song_id INT NOT NULL,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (library_id, song_id),
    FOREIGN KEY (library_id) REFERENCES libraries(id),
    FOREIGN KEY (song_id) REFERENCES songs(id)
);

-- Tabla de reproducciones
CREATE TABLE playbacks (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    song_id INT NOT NULL,
    played_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('playing', 'paused', 'completed') NOT NULL,
    duration INT NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (song_id) REFERENCES songs(id)
);

-- Insertar usuarios admin por defecto
INSERT INTO users (name, email, password, role) VALUES
('HENRY ALIAGA', 'henry@example.com', 'admin123', 'admin'),
('ISMAEL ESPINOZA', 'ismael@example.com', 'admin123', 'admin');

-- Insertar usuarios simulados
INSERT INTO users (name, email, password, role) VALUES
('Juan Perez', 'juan.perez@example.com', 'password123', 'user'),
('Ana Gomez', 'ana.gomez@example.com', 'securepass456', 'user'),
('Carlos Lopez', 'carlos.lopez@example.com', 'qwerty789', 'user'); 

