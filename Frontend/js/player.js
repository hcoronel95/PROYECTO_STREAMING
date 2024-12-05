// Frontend/js/player.js
class MusicPlayer {
    constructor() {
        this.currentSong = null;
        this.isPlaying = false;
        this.songList = [];
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
            }
        } catch (error) {
            console.error('Error cargando canciones:', error);
        }
    }

    renderSongList() {
        this.songListElement.innerHTML = this.songList.map(song => `
            <div class="list-group-item song-list-item ${this.currentSong?.id === song.id ? 'active' : ''}"
                 data-id="${song.id}">
                <div class="d-flex justify-content-between align-items-center">
                    <div>
                        <h6 class="mb-0">${song.title}</h6>
                        <small class="text-muted">${song.artist}</small>
                    </div>
                    <span class="badge bg-secondary">${song.genre}</span>
                </div>
            </div>
        `).join('');
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
                this.updatePlayerUI();
            }
        } catch (error) {
            console.error('Error reproduciendo canción:', error);
        }
    }

    togglePlay() {
        if (!this.currentSong) return;
        
        this.isPlaying = !this.isPlaying;
        this.updatePlayerUI();
        
        fetch(`/api/${this.isPlaying ? 'play' : 'pause'}/${this.currentSong.id}`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('userToken')}`
            }
        });
    }

    playPrevious() {
        if (!this.currentSong) return;
        const currentIndex = this.songList.findIndex(song => song.id === this.currentSong.id);
        if (currentIndex > 0) {
            this.playSong(this.songList[currentIndex - 1].id);
        }
    }

    playNext() {
        if (!this.currentSong) return;
        const currentIndex = this.songList.findIndex(song => song.id === this.currentSong.id);
        if (currentIndex < this.songList.length - 1) {
            this.playSong(this.songList[currentIndex + 1].id);
        }
    }

    updatePlayerUI() {
        this.playBtn.textContent = this.isPlaying ? '⏸' : '▶';
        if (this.currentSong) {
            this.currentSongElement.textContent = 
                `${this.currentSong.title} - ${this.currentSong.artist}`;
        }
        this.renderSongList();
    }
}

// Inicializar el reproductor cuando se carga la página
document.addEventListener('DOMContentLoaded', () => {
    new MusicPlayer();
});