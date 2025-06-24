document.addEventListener('DOMContentLoaded', () => {
    console.log('Frontend JavaScript is running!');

    const movieGrid = document.querySelector('.trending-movies .movie-grid');

    async function fetchTrendingMovies() {
        try {
            const response = await fetch('/api/trending');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const movies = await response.json();
            displayMovies(movies);
        } catch (error) {
            console.error('Error fetching trending movies:', error);
            movieGrid.innerHTML = '<p>Failed to load trending movies. Please try again later.</p>';
        }
    }

    function displayMovies(movies) {
        movieGrid.innerHTML = ''; // Clear existing content
        if (movies && movies.length > 0) {
            movies.forEach(movie => {
                const movieCard = document.createElement('div');
                movieCard.classList.add('movie-card');

                const posterPath = movie.poster_path ? `https://image.tmdb.org/t/p/w200${movie.poster_path}` : 'https://via.placeholder.com/200x300?text=No+Image';
                const title = movie.title || movie.name || 'Untitled';
                const releaseDate = movie.release_date ? new Date(movie.release_date).getFullYear() : 'N/A';
                const voteAverage = movie.vote_average ? movie.vote_average.toFixed(1) : 'N/A';
                const genres = movie.genre_ids ? movie.genre_ids.map(id => getGenreName(id)).join(', ') : 'N/A';

                movieCard.innerHTML = `
                    <img src="${posterPath}" alt="${title}">
                    <div class="movie-info">
                        <h3>${title}</h3>
                        <p>${releaseDate} &bull; ${genres}</p>
                        <div class="rating"><i class="fas fa-star"></i> ${voteAverage}</div>
                    </div>
                `;
                movieGrid.appendChild(movieCard);
            });
        } else {
            movieGrid.innerHTML = '<p>No trending movies found.</p>';
        }
    }

    // Simple genre mapping (you might want a more comprehensive list)
    function getGenreName(genreId) {
        const genres = {
            28: 'Action',
            12: 'Adventure',
            16: 'Animation',
            35: 'Comedy',
            80: 'Crime',
            99: 'Documentary',
            18: 'Drama',
            10751: 'Family',
            14: 'Fantasy',
            36: 'History',
            27: 'Horror',
            10402: 'Music',
            9648: 'Mystery',
            10749: 'Romance',
            878: 'Science Fiction',
            10770: 'TV Movie',
            53: 'Thriller',
            10752: 'War',
            37: 'Western'
        };
        return genres[genreId] || 'Unknown';
    }

    // Initial fetch when the page loads
    fetchTrendingMovies();
});