let player = null; // Variable global para el reproductor

class MusicPlayer {
    constructor() {
        this.currentSong = null;
        this.audio = new Audio();
        this.songs = [];
        this.isPlaying = false;
        this.initializePlayer();
        this.loadSongs();
    }

    initializePlayer() {
        const playBtn = document.querySelector('.player-btn.play');
        const prevBtn = document.querySelector('.player-btn.prev');
        const nextBtn = document.querySelector('.player-btn.next');
        const searchInput = document.querySelector('.search-input');
        const progressBar = document.querySelector('.progress-bar');

        if (playBtn) {
            playBtn.addEventListener('click', () => this.togglePlay());
        }
        if (prevBtn) {
            prevBtn.addEventListener('click', () => this.playPrevious());
        }
        if (nextBtn) {
            nextBtn.addEventListener('click', () => this.playNext());
        }
        if (searchInput) {
            searchInput.addEventListener('input', (e) => this.searchSongs(e.target.value));
        }
        if (progressBar) {
            progressBar.addEventListener('click', (e) => this.setProgress(e));
        }

        this.audio.addEventListener('timeupdate', () => this.updateProgress());
        this.audio.addEventListener('ended', () => this.playNext());
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
            } else {
                console.error('Error al cargar canciones:', await response.text());
            }
        } catch (error) {
            console.error('Error cargando canciones:', error);
        }
    }

    displaySongs(songs) {
        const musicList = document.querySelector('.music-list');
        if (!musicList) return;

        musicList.innerHTML = songs.map((song, index) => `
            <div class="song-item" data-index="${index}">
                <img src="../images/music-cover.jpg" alt="${song.title}" class="custom-img-size">
                <div class="song-info">
                    <h3>${song.title}</h3>
                    <p>${song.artist}</p>
                    <p>${song.genre}</p>
                </div>
            </div>
        `).join('');

        musicList.querySelectorAll('.song-item').forEach(item => {
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

        const titleElement = document.querySelector('.song-title');
        const artistElement = document.querySelector('.song-artist');

        if (titleElement) titleElement.textContent = song.title;
        if (artistElement) artistElement.textContent = song.artist;

        this.audio.src = song.file_path;
        this.audio.play()
            .catch(error => console.error('Error reproduciendo canción:', error));
        this.isPlaying = true;
        this.updatePlayButton();
    }

    togglePlay() {
        if (this.audio.paused) {
            this.audio.play().catch(error => console.error('Error reproduciendo:', error));
            this.isPlaying = true;
        } else {
            this.audio.pause();
            this.isPlaying = false;
        }
        this.updatePlayButton();
    }

    updatePlayButton() {
        const playBtn = document.querySelector('.player-btn.play');
        if (playBtn) {
            playBtn.textContent = this.isPlaying ? '⏸' : '▶';
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

    updateProgress() {
        const progress = document.querySelector('.progress');
        const progressBar = document.querySelector('.progress-bar');
        if (progress && progressBar && this.audio.duration) {
            const percentage = (this.audio.currentTime / this.audio.duration) * 100;
            progress.style.width = `${percentage}%`;
        }
    }

    setProgress(e) {
        const progressBar = document.querySelector('.progress-bar');
        if (progressBar && this.audio.duration) {
            const width = progressBar.clientWidth;
            const clickX = e.offsetX;
            const duration = this.audio.duration;
            this.audio.currentTime = (clickX / width) * duration;
        }
    }

    searchSongs(query) {
        const filteredSongs = this.songs.filter(song =>
            song.title.toLowerCase().includes(query.toLowerCase()) ||
            song.artist.toLowerCase().includes(query.toLowerCase()) ||
            song.genre.toLowerCase().includes(query.toLowerCase())
        );
        this.displaySongs(filteredSongs);
    }
}

document.addEventListener("DOMContentLoaded", () => {
    player = new MusicPlayer();

    const token = localStorage.getItem('userToken');
    if (token) {
        fetch('/api/user-info', {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        })
            .then(response => response.json())
            .then(data => {
                const welcomeMessage = document.querySelector('#welcome-message');
                if (data && data.name) {
                    welcomeMessage.textContent = `Bienvenido, ${data.name}`;
                } else {
                    welcomeMessage.textContent = 'Bienvenido, usuario desconocido';
                }
            })
            .catch(error => console.error('Error obteniendo información del usuario:', error));
    }

    const logoutButton = document.querySelector(".logout-button");
    if (logoutButton) {
        logoutButton.addEventListener("click", () => {
            if (confirm("¿Estás seguro de que deseas cerrar sesión?")) {
                fetch("/api/logout", { method: "POST", headers: { "Content-Type": "application/json" } })
                    .then(response => {
                        if (response.ok) {
                            alert("Sesión cerrada exitosamente.");
                            window.location.href = "/";
                        } else {
                            alert("Error cerrando sesión.");
                        }
                    })
                    .catch(err => console.error("Error:", err));
            }
        });
    }
});
