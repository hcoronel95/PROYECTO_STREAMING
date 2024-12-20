class MusicPlayer {
    constructor() {
        this.currentSong = null;
        this.audio = new Audio();
        this.songs = [];
        this.isPlaying = false;
        this.albumCoverElement = document.getElementById('albumCover');
        this.initializeElements();
        this.loadSongs();
        this.setupEventListeners();
    }

    initializeElements() {
        this.playBtn = document.getElementById('playBtn');
        this.prevBtn = document.getElementById('prevBtn');
        this.nextBtn = document.getElementById('nextBtn');
        this.songListElement = document.getElementById('songList');
        this.currentSongElement = document.getElementById('currentSong');
        this.currentArtistElement = document.getElementById('currentArtist');
    }

    async loadSongs() {
        try {
            const response = await fetch('/api/songs/list', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('userToken')}`
                }
            });
            
            if (response.ok) {
                this.songs = await response.json();
                console.log('Canciones cargadas:', this.songs);
                this.displaySongs(this.songs);
            } else {
                console.error('Error al cargar canciones:', response.statusText);
            }
        } catch (error) {
            console.error('Error cargando canciones:', error);
        }
    }

    displaySongs(songs) {
        const musicList = document.querySelector('.song-list-container');
        if (!musicList) return;

        const songItems = songs.map((song, index) => `
            <div class="song-item" data-index="${index}">
                <img src="../images/music-cover.jpg" alt="${song.title}" class="custom-img-size">
                <div class="song-info">
                    <h3>${song.title}</h3>
                    <p>${song.artist}</p>
                    <p>${song.genre}</p>
                </div>
            </div>
        `).join('');

        musicList.innerHTML = songItems;

        // Agregar event listeners a las canciones
        document.querySelectorAll('.song-item').forEach(item => {
            item.addEventListener('click', () => {
                const index = parseInt(item.dataset.index);
                this.playSong(index);
            });
        });
    }

    playSong(index) {
        if (index < 0 || index >= this.songs.length) return;
        
        const song = this.songs[index];
        this.currentSong = index;
        
        // Construir la URL completa para el archivo de audio
        const fileName = song.file_path.split('\\').pop(); // Manejar rutas de Windows
        const audioUrl = `/uploads/songs/${fileName}`;
        console.log('Intentando reproducir:', audioUrl);
        
        // Actualizar interfaz
        if (this.currentSongElement) this.currentSongElement.textContent = song.title;
        if (this.currentArtistElement) this.currentArtistElement.textContent = song.artist;
        if (this.albumCoverElement) this.albumCoverElement.src = '../images/music-cover.jpg';
        
        this.audio.src = audioUrl;
        this.audio.play()
            .catch(error => {
                console.error('Error reproduciendo canción:', error);
            });
        
        this.isPlaying = true;
        this.updatePlayButton();
    }

    togglePlay() {
        if (this.audio.paused) {
            this.audio.play()
                .catch(error => {
                    console.error('Error reproduciendo:', error);
                });
            this.isPlaying = true;
        } else {
            this.audio.pause();
            this.isPlaying = false;
        }
        this.updatePlayButton();
    }

    updatePlayButton() {
        if (this.playBtn) {
            this.playBtn.textContent = this.isPlaying ? '⏸' : '▶';
        }
    }

    playNext() {
        if (this.currentSong !== null) {
            const nextIndex = (this.currentSong + 1) % this.songs.length;
            this.playSong(nextIndex);
        }
    }

    playPrevious() {
        if (this.currentSong !== null) {
            const prevIndex = (this.currentSong - 1 + this.songs.length) % this.songs.length;
            this.playSong(prevIndex);
        }
    }

    setupEventListeners() {
        if (this.playBtn) {
            this.playBtn.addEventListener('click', () => this.togglePlay());
        }
        if (this.prevBtn) {
            this.prevBtn.addEventListener('click', () => this.playPrevious());
        }
        if (this.nextBtn) {
            this.nextBtn.addEventListener('click', () => this.playNext());
        }

        // Evento para cuando termine la canción
        this.audio.addEventListener('ended', () => this.playNext());
    }
}

// Inicializar el reproductor cuando se cargue el DOM
document.addEventListener('DOMContentLoaded', () => {
    const player = new MusicPlayer();
});

// Evento para el botón de logout
document.querySelector(".logout-btn")?.addEventListener("click", () => {
    if (confirm("¿Estás seguro de que deseas cerrar sesión?")) {
        fetch("/api/logout", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
        })
            .then((response) => {
                if (response.ok) {
                    localStorage.removeItem('userToken');
                    window.location.href = "/";
                } else {
                    alert("Error cerrando sesión.");
                }
            })
            .catch((err) => console.error("Error:", err));
    }
});