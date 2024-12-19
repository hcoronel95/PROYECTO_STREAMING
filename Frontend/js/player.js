class MusicPlayer {
    constructor() {
        this.currentSong = null;
        this.isPlaying = false;
        this.songList = [];
        this.favorites = []; // Almacena las canciones favoritas del usuario
        this.albumCoverElement = document.getElementById('albumCover');
        this.initializeElements();
        this.loadSongs();
        this.loadFavorites(); // Carga las canciones favoritas
        this.setupEventListeners();
    }

    initializeElements() {
        this.playBtn = document.getElementById('playBtn');
        this.prevBtn = document.getElementById('prevBtn');
        this.nextBtn = document.getElementById('nextBtn');
        this.songListElement = document.getElementById('songList');
        this.favoritesListElement = document.getElementById('favoritesList'); // Nuevo para mostrar favoritos
        this.currentSongElement = document.getElementById('currentSong');
    }

    async loadSongs() {
        try {
            const response = await fetch('/api/songs', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('userToken')}`
                }
            });
            if (response.ok) {
                this.songList = await response.json();
                this.renderSongList();
            } else {
                console.error('Error al cargar canciones:', response.statusText);
            }
        } catch (error) {
            console.error('Error cargando canciones:', error);
        }
    }

    async loadFavorites() {
        try {
            const response = await fetch('/api/favorites', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('userToken')}`
                }
            });
            if (response.ok) {
                this.favorites = await response.json();
                this.renderFavorites();
            } else {
                console.error('Error al cargar favoritos:', response.statusText);
            }
        } catch (error) {
            console.error('Error cargando favoritos:', error);
        }
    }

    renderSongList() {
        this.songListElement.innerHTML = this.songList.map(song => 
            `<div class="list-group-item song-list-item ${this.currentSong?.id === song.id ? 'active' : ''}"
                 data-id="${song.id}">
                <div class="d-flex justify-content-between align-items-center">
                    <div>
                        <h6 class="mb-0">${song.title}</h6>
                        <small class="text-muted">${song.artist}</small>
                    </div>
                    <button class="btn btn-sm btn-outline-secondary add-favorite-btn" data-id="${song.id}">
                        ${this.favorites.some(fav => fav.id === song.id) ? '❤️' : '🤍'}
                    </button>
                </div>
            </div>`).join('');
    }

    renderFavorites() {
        this.favoritesListElement.innerHTML = this.favorites.map(song =>
            `<div class="list-group-item">
                <h6 class="mb-0">${song.title}</h6>
                <small class="text-muted">${song.artist}</small>
            </div>`).join('');
    }

    setupEventListeners() {
        this.playBtn.addEventListener('click', () => this.togglePlay());
        this.prevBtn.addEventListener('click', () => this.playPrevious());
        this.nextBtn.addEventListener('click', () => this.playNext());
        
        this.songListElement.addEventListener('click', (e) => {
            const songItem = e.target.closest('.song-list-item');
            if (songItem) {
                const songId = parseInt(songItem.dataset.id);
                this.playSong(songId);
            }

            const favoriteBtn = e.target.closest('.add-favorite-btn');
            if (favoriteBtn) {
                const songId = parseInt(favoriteBtn.dataset.id);
                this.toggleFavorite(songId);
            }
        });
    }

    async playSong(id) {
        try {
            const response = await fetch(`/api/songs/play/${id}`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('userToken')}`
                }
            });

            if (response.ok) {
                this.currentSong = this.songList.find(song => song.id === id);
                this.isPlaying = true;
                this.albumCoverElement.src = `/path/to/song-${this.currentSong.id}-cover.jpg`;
                this.updatePlayerUI();
            } else {
                console.error('Error al reproducir canción:', response.statusText);
            }
        } catch (error) {
            console.error('Error reproduciendo canción:', error);
        }
    }

    async toggleFavorite(id) {
        const isFavorite = this.favorites.some(song => song.id === id);

        try {
            const response = await fetch(`/api/favorites/${isFavorite ? 'remove' : 'add'}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('userToken')}`
                },
                body: JSON.stringify({ songId: id })
            });

            if (response.ok) {
                this.favorites = isFavorite
                    ? this.favorites.filter(song => song.id !== id)
                    : [...this.favorites, this.songList.find(song => song.id === id)];
                this.renderSongList();
                this.renderFavorites();
            } else {
                console.error('Error al actualizar favoritos:', response.statusText);
            }
        } catch (error) {
            console.error('Error actualizando favoritos:', error);
        }
    }

    togglePlay() {
        if (!this.currentSong) return;
        this.isPlaying = !this.isPlaying;
        this.updatePlayerUI();
    }

    playPrevious() {
        const index = this.songList.findIndex(song => song.id === this.currentSong.id);
        if (index > 0) this.playSong(this.songList[index - 1].id);
    }

    playNext() {
        const index = this.songList.findIndex(song => song.id === this.currentSong.id);
        if (index < this.songList.length - 1) this.playSong(this.songList[index + 1].id);
    }

    updatePlayerUI() {
        this.playBtn.textContent = this.isPlaying ? '⏸' : '▶';
        this.currentSongElement.textContent = this.currentSong ? `${this.currentSong.title} - ${this.currentSong.artist}` : 'No hay canción seleccionada';
        this.renderSongList();
    }
}

// Inicializar el MusicPlayer al cargar el DOM
document.addEventListener('DOMContentLoaded', () => {
    new MusicPlayer();
});

// Evento para el botón de logout
document.querySelector(".logout-btn").addEventListener("click", () => {
    if (confirm("¿Estás seguro de que deseas cerrar sesión?")) {
        fetch("/api/logout", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
        })
            .then((response) => {
                if (response.ok) {
                    alert("Sesión cerrada exitosamente.");
                    window.location.href = "/"; // Redirige al inicio
                } else {
                    alert("Error cerrando sesión.");
                }
            })
            .catch((err) => console.error("Error:", err));
    }
});
