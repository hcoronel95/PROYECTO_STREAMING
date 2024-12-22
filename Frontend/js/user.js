class MusicPlayer {
    constructor() {
        this.currentSong = null;
        this.audio = new Audio();
        this.songs = [];
        this.isPlaying = false;
        this.currentSongIndex = -1;
        this.initializeElements();
        this.loadUserInfo();
        this.loadSongs();
        this.setupEventListeners();
    }

    initializeElements() {
        this.songGrid = document.getElementById('songGrid');
        this.searchInput = document.querySelector('.search-input');
        this.genreFilter = document.querySelector('.genre-filter');
        this.playBtn = document.getElementById('playBtn');
        this.prevBtn = document.getElementById('prevBtn');
        this.nextBtn = document.getElementById('nextBtn');
        this.currentTitle = document.getElementById('currentTitle');
        this.currentArtist = document.getElementById('currentArtist');
        this.currentCover = document.getElementById('currentCover');
        this.progress = document.querySelector('.progress');
    }

    setupEventListeners() {
        if (this.searchInput) {
            this.searchInput.addEventListener('input', (e) => this.filterSongs(e.target.value));
        }
        if (this.genreFilter) {
            this.genreFilter.addEventListener('change', (e) => this.filterByGenre(e.target.value));
        }
        if (this.playBtn) {
            this.playBtn.addEventListener('click', () => this.togglePlay());
        }
        if (this.prevBtn) {
            this.prevBtn.addEventListener('click', () => this.playPrevious());
        }
        if (this.nextBtn) {
            this.nextBtn.addEventListener('click', () => this.playNext());
        }

        // Eventos del reproductor de audio
        this.audio.addEventListener('timeupdate', () => this.updateProgress());
        this.audio.addEventListener('ended', () => this.playNext());

        // Evento de logout
        const logoutButton = document.querySelector('.logout-button');
        if (logoutButton) {
            logoutButton.addEventListener('click', () => this.handleLogout());
        }
    }

    async loadUserInfo() {
        try {
            const response = await fetch('/api/user-info', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('userToken')}`
                }
            });
            if (response.ok) {
                const userData = await response.json();
                const welcomeMessage = document.getElementById('welcome-message');
                if (welcomeMessage && userData.name) {
                    welcomeMessage.textContent = `Bienvenido`;
                }
            }
        } catch (error) {
            console.error('Error cargando información del usuario:', error);
        }
    }

    async loadSongs() {
        try {
            const response = await fetch('/api/songs', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('userToken')}`
                }
            });
            if (response.ok) {
                this.songs = await response.json();
                this.displaySongs(this.songs);
            }
        } catch (error) {
            console.error('Error cargando canciones:', error);
        }
    }

    displaySongs(songs) {
        if (!this.songGrid) return;

        this.songGrid.innerHTML = songs.map((song, index) => `
            <div class="song-card" data-index="${index}">
                <img src="../images/${song.genre.toLowerCase()}.jpg" 
                     alt="${song.title}" 
                     class="song-image"
                     onerror="this.src='../images/fondonota.jpg'">
                <div class="song-info">
                    <h3>${song.title}</h3>
                    <p>${song.artist}</p>
                    <p class="genre-tag">${song.genre}</p>
                </div>
            </div>
        `).join('');

        // Agregar event listeners a las tarjetas
        this.songGrid.querySelectorAll('.song-card').forEach(card => {
            card.addEventListener('click', () => {
                const index = parseInt(card.dataset.index);
                this.playSong(index);
            });

            // Preview al pasar el mouse
            card.addEventListener('mouseenter', () => {
                if (!this.isPlaying) {
                    const index = parseInt(card.dataset.index);
                    this.previewSong(index);
                }
            });

            card.addEventListener('mouseleave', () => {
                if (!this.isPlaying) {
                    this.stopPreview();
                }
            });
        });
    }

    playSong(index) {
        if (index < 0 || index >= this.songs.length) return;

        const song = this.songs[index];
        this.currentSongIndex = index;
        this.currentSong = song;

        // Actualizar interfaz
        this.currentTitle.textContent = song.title;
        this.currentArtist.textContent = song.artist;
        this.currentCover.src = `../images/${song.genre.toLowerCase()}.jpg`;

        // Reproducir audio
        this.audio.src = song.file_path;
        this.audio.play()
            .then(() => {
                this.isPlaying = true;
                this.updatePlayButton();
            })
            .catch(error => console.error('Error reproduciendo:', error));
    }

    previewSong(index) {
        const song = this.songs[index];
        this.audio.src = song.file_path;
        this.audio.volume = 0.3; // Volumen más bajo para preview
        this.audio.currentTime = 0;
        this.audio.play().catch(error => console.error('Error en preview:', error));
    }

    stopPreview() {
        this.audio.pause();
        this.audio.currentTime = 0;
        this.audio.volume = 1.0; // Restaurar volumen normal
    }

    togglePlay() {
        if (this.currentSongIndex === -1 && this.songs.length > 0) {
            this.playSong(0);
        } else if (this.audio.paused) {
            this.audio.play();
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
        if (this.currentSongIndex > -1) {
            const nextIndex = (this.currentSongIndex + 1) % this.songs.length;
            this.playSong(nextIndex);
        }
    }

    playPrevious() {
        if (this.currentSongIndex > -1) {
            const prevIndex = (this.currentSongIndex - 1 + this.songs.length) % this.songs.length;
            this.playSong(prevIndex);
        }
    }

    updateProgress() {
        if (this.progress && this.audio.duration) {
            const percentage = (this.audio.currentTime / this.audio.duration) * 100;
            this.progress.style.width = `${percentage}%`;
        }
    }

    filterSongs(query) {
        const filteredSongs = this.songs.filter(song =>
            song.title.toLowerCase().includes(query.toLowerCase()) ||
            song.artist.toLowerCase().includes(query.toLowerCase()) ||
            song.genre.toLowerCase().includes(query.toLowerCase())
        );
        this.displaySongs(filteredSongs);
    }

    filterByGenre(genre) {
        if (!genre) {
            this.displaySongs(this.songs);
        } else {
            const filteredSongs = this.songs.filter(song => 
                song.genre.toLowerCase() === genre.toLowerCase()
            );
            this.displaySongs(filteredSongs);
        }
    }

    async handleLogout() {
        if (confirm('¿Estás seguro de que deseas cerrar sesión?')) {
            try {
                const response = await fetch('/api/logout', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' }
                });

                if (response.ok) {
                    localStorage.removeItem('userToken');
                    window.location.href = '/';
                } else {
                    throw new Error('Error en logout');
                }
            } catch (error) {
                console.error('Error:', error);
                alert('Error al cerrar sesión');
            }
        }
    }
}

// Inicializar el reproductor cuando el DOM esté listo
document.addEventListener('DOMContentLoaded', () => {
    new MusicPlayer();
});