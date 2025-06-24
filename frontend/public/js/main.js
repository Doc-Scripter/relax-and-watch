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

                const posterPath = movie.poster_path ? `https://image.tmdb.org/t/p/w200${movie.poster_path}` : 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjMwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjMzMzIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIxNCIgZmlsbD0iI2ZmZiIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPk5vIEltYWdlPC90ZXh0Pjwvc3ZnPg==';
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
                movieCard.dataset.movieId = movie.id; // Store movie ID
                movieCard.addEventListener('click', () => fetchMovieDetails(movie.id));
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

    async function fetchMovieDetails(movieId) {
        try {
            const response = await fetch(`/api/movie/${movieId}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const movieDetails = await response.json();
            displayMovieDetailsModal(movieDetails);
        } catch (error) {
            console.error('Error fetching movie details:', error);
            alert('Failed to load movie details. Please try again later.');
        }
    }

    function displayMovieDetailsModal(movieData) {
        console.log('Movie Details:', movieData);
        const modal = document.getElementById('movieDetailsModal');
        const modalBody = document.getElementById('modal-body');

        // Extract data from both TMDB and OMDB sources
        const tmdbData = movieData.TMDBData || {};
        const omdbData = movieData.OMDBData || {};

        // Use OMDB data as primary source for detailed info, fallback to TMDB
        const title = omdbData.Title || tmdbData.title || 'Unknown Title';
        const year = omdbData.Year || (tmdbData.release_date ? new Date(tmdbData.release_date).getFullYear() : 'N/A');
        const rated = omdbData.Rated || 'N/A';
        const released = omdbData.Released || tmdbData.release_date || 'N/A';
        const runtime = omdbData.Runtime || (tmdbData.runtime ? `${tmdbData.runtime} min` : 'N/A');
        const genre = omdbData.Genre || (tmdbData.genres ? tmdbData.genres.map(g => g.name).join(', ') : 'N/A');
        const director = omdbData.Director || 'N/A';
        const writer = omdbData.Writer || 'N/A';
        const actors = omdbData.Actors || 'N/A';
        const plot = omdbData.Plot || tmdbData.overview || 'No plot available';
        const language = omdbData.Language || (tmdbData.original_language || 'N/A');
        const country = omdbData.Country || (tmdbData.production_countries ? tmdbData.production_countries.map(c => c.name).join(', ') : 'N/A');
        const awards = omdbData.Awards || 'N/A';
        const imdbRating = omdbData.imdbRating || (tmdbData.vote_average ? tmdbData.vote_average.toFixed(1) : 'N/A');
        const boxOffice = omdbData.BoxOffice || 'N/A';
        const posterPath = tmdbData.poster_path ? `https://image.tmdb.org/t/p/w300${tmdbData.poster_path}` : 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzAwIiBoZWlnaHQ9IjQ1MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjMzMzIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIxOCIgZmlsbD0iI2ZmZiIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPk5vIEltYWdlPC90ZXh0Pjwvc3ZnPg==';

        modalBody.innerHTML = `
            <div class="movie-details-container">
                <div class="movie-poster">
                    <img src="${posterPath}" alt="${title}" style="max-width: 300px; border-radius: 8px;">
                </div>
                <div class="movie-info-detailed">
                    <h2>${title}</h2>
                    <p><strong>Year:</strong> ${year}</p>
                    <p><strong>Rated:</strong> ${rated}</p>
                    <p><strong>Released:</strong> ${released}</p>
                    <p><strong>Runtime:</strong> ${runtime}</p>
                    <p><strong>Genre:</strong> ${genre}</p>
                    <p><strong>Director:</strong> ${director}</p>
                    <p><strong>Writer:</strong> ${writer}</p>
                    <p><strong>Actors:</strong> ${actors}</p>
                    <p><strong>Plot:</strong> ${plot}</p>
                    <p><strong>Language:</strong> ${language}</p>
                    <p><strong>Country:</strong> ${country}</p>
                    <p><strong>Awards:</strong> ${awards}</p>
                    <p><strong>IMDB Rating:</strong> ${imdbRating}</p>
                    <p><strong>Box Office:</strong> ${boxOffice}</p>
                </div>
            </div>
        `;
        modal.style.display = 'block';

        const closeButton = document.querySelector('.close-button');
        closeButton.onclick = function() {
            modal.style.display = 'none';
        }

        window.onclick = function(event) {
            if (event.target == modal) {
                modal.style.display = 'none';
            }
        }
    }

    // Initial fetch when the page loads
    fetchTrendingMovies();
});