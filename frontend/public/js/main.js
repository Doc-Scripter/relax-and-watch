document.addEventListener('DOMContentLoaded', () => {
    console.log('Frontend JavaScript is running!');

    const movieGrid = document.querySelector('.trending-movies .movie-grid');
    const searchInput = document.querySelector('.search-bar input');
    const searchButton = document.querySelector('.search-bar button');
    const movieButton = document.querySelector('.filter-buttons .btn:first-child');
    const seriesButton = document.querySelector('.filter-buttons .btn:last-child');
    const watchlistButton = document.querySelector('.nav-links .btn:last-child');
    const genreSelect = document.getElementById('genre');
    const yearSelect = document.getElementById('year');
    const ratingSelect = document.getElementById('rating');
    const runtimeSelect = document.getElementById('runtime');
    const trendingSection = document.querySelector('.trending-movies h2');
    
    // State management
    let currentContentType = 'movie'; // 'movie' or 'tv'
    let currentPage = 1;
    let totalPages = 1;
    let isLoading = false;
    let searchTimeout;
    let watchlist = JSON.parse(localStorage.getItem('watchlist')) || [];
    let currentQuery = '';
    let currentFilters = {};

    // Initialize the app
    init();

    async function init() {
        updateWatchlistCount();
        await loadGenres();
        await loadYears();
        await loadRatings();
        await loadRuntimes();
        await fetchTrendingMovies();
        setupEventListeners();
    }

    function setupEventListeners() {
        // Search functionality with debounce
        searchInput.addEventListener('input', handleSearchInput);
        searchButton.addEventListener('click', handleSearchClick);
        
        // Content type buttons
        movieButton.addEventListener('click', () => switchContentType('movie'));
        seriesButton.addEventListener('click', () => switchContentType('tv'));
        
        // Watchlist button
        watchlistButton.addEventListener('click', showWatchlist);
        
        // Filter dropdowns
        genreSelect.addEventListener('change', handleFilterChange);
        yearSelect.addEventListener('change', handleFilterChange);
        ratingSelect.addEventListener('change', handleFilterChange);
        runtimeSelect.addEventListener('change', handleFilterChange);
        
        // Infinite scroll for pagination
        window.addEventListener('scroll', handleScroll);
    }

    function handleSearchInput(e) {
        const query = e.target.value.trim();
        
        // Clear previous timeout
        if (searchTimeout) {
            clearTimeout(searchTimeout);
        }
        
        // Set new timeout for debounce (500ms delay)
        searchTimeout = setTimeout(() => {
            if (query.length > 2) {
                currentQuery = query;
                currentPage = 1;
                searchContent(query);
            } else if (query.length === 0) {
                currentQuery = '';
                currentPage = 1;
                fetchTrendingMovies();
            }
        }, 500);
    }

    function handleSearchClick() {
        const query = searchInput.value.trim();
        if (query.length > 2) {
            currentQuery = query;
            currentPage = 1;
            searchContent(query);
        }
    }

    async function switchContentType(type) {
        currentContentType = type;
        currentPage = 1;
        
        // Update button states
        movieButton.classList.toggle('active', type === 'movie');
        seriesButton.classList.toggle('active', type === 'tv');
        
        // Update section title
        trendingSection.textContent = type === 'movie' ? 'Trending Movies' : 'Trending TV Series';
        
        // Reload genres for the new content type
        await loadGenres();
        
        // Reset filters
        currentFilters = {};
        genreSelect.value = 'all';
        yearSelect.value = 'all';
        ratingSelect.value = 'all';
        runtimeSelect.value = 'all';
        
        // Fetch new content
        if (currentQuery) {
            searchContent(currentQuery);
        } else {
            fetchTrendingMovies();
        }
    }

    async function fetchTrendingMovies() {
        if (isLoading) return;
        isLoading = true;
        
        try {
            let url = '/api/trending';
            const params = new URLSearchParams();
            
            if (currentContentType === 'tv') {
                params.append('type', 'tv');
            }
            
            if (currentPage > 1) {
                params.append('page', currentPage.toString());
            }
            
            if (params.toString()) {
                url += '?' + params.toString();
            }
            
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            
            if (currentPage === 1) {
                displayMovies(data.results || data);
            } else {
                appendMovies(data.results || data);
            }
            
            totalPages = data.total_pages || 1;
        } catch (error) {
            console.error('Error fetching trending content:', error);
            if (currentPage === 1) {
                movieGrid.innerHTML = '<p>Failed to load trending content. Please try again later.</p>';
            }
        } finally {
            isLoading = false;
        }
    }

    async function searchContent(query) {
        if (isLoading) return;
        isLoading = true;
        
        try {
            const params = new URLSearchParams();
            params.append('q', query);
            params.append('type', currentContentType);
            
            if (currentPage > 1) {
                params.append('page', currentPage.toString());
            }
            
            const response = await fetch(`/api/search?${params.toString()}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            
            if (currentPage === 1) {
                displayMovies(data.results || data);
            } else {
                appendMovies(data.results || data);
            }
            
            totalPages = data.total_pages || 1;
        } catch (error) {
            console.error('Error searching content:', error);
            if (currentPage === 1) {
                movieGrid.innerHTML = '<p>Failed to search content. Please try again later.</p>';
            }
        } finally {
            isLoading = false;
        }
    }

    async function discoverContent() {
        if (isLoading) return;
        isLoading = true;
        
        try {
            const params = new URLSearchParams();
            params.append('type', currentContentType);
            
            if (currentPage > 1) {
                params.append('page', currentPage.toString());
            }
            
            // Add filters
            Object.keys(currentFilters).forEach(key => {
                if (currentFilters[key]) {
                    params.append(key, currentFilters[key]);
                }
            });
            
            const response = await fetch(`/api/discover?${params.toString()}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            
            if (currentPage === 1) {
                displayMovies(data.results || data);
            } else {
                appendMovies(data.results || data);
            }
            
            totalPages = data.total_pages || 1;
        } catch (error) {
            console.error('Error discovering content:', error);
            if (currentPage === 1) {
                movieGrid.innerHTML = '<p>Failed to discover content. Please try again later.</p>';
            }
        } finally {
            isLoading = false;
        }
    }

    function handleFilterChange() {
        currentFilters = {
            genre: genreSelect.value !== 'all' ? genreSelect.value : '',
            year: yearSelect.value !== 'all' ? yearSelect.value : '',
            rating: ratingSelect.value !== 'all' ? ratingSelect.value : '',
            runtime: runtimeSelect.value !== 'all' ? runtimeSelect.value : ''
        };
        
        // Remove empty filters
        Object.keys(currentFilters).forEach(key => {
            if (!currentFilters[key]) {
                delete currentFilters[key];
            }
        });
        
        currentPage = 1;
        currentQuery = '';
        searchInput.value = '';
        
        if (Object.keys(currentFilters).length > 0) {
            discoverContent();
        } else {
            fetchTrendingMovies();
        }
    }

    function handleScroll() {
        if (isLoading || currentPage >= totalPages) return;
        
        const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        const windowHeight = window.innerHeight;
        const documentHeight = document.documentElement.scrollHeight;
        
        // Load more when user is 200px from bottom
        if (scrollTop + windowHeight >= documentHeight - 200) {
            currentPage++;
            
            if (currentQuery) {
                searchContent(currentQuery);
            } else if (Object.keys(currentFilters).length > 0) {
                discoverContent();
            } else {
                fetchTrendingMovies();
            }
        }
    }

    async function loadGenres() {
        try {
            const response = await fetch(`/api/genres?type=${currentContentType}`);
            if (response.ok) {
                const data = await response.json();
                const genres = data.genres || data || [];
                
                genreSelect.innerHTML = '<option value="all">All Genres</option>';
                genres.forEach(genre => {
                    const option = document.createElement('option');
                    option.value = genre.id;
                    option.textContent = genre.name;
                    genreSelect.appendChild(option);
                });
            }
        } catch (error) {
            console.error('Error loading genres:', error);
        }
    }

    async function loadYears() {
        const currentYear = new Date().getFullYear();
        yearSelect.innerHTML = '<option value="all">All Years</option>';
        
        for (let year = currentYear; year >= 1900; year--) {
            const option = document.createElement('option');
            option.value = year;
            option.textContent = year;
            yearSelect.appendChild(option);
        }
    }

    async function loadRatings() {
        const ratings = [
            { value: '9', text: '9.0+ Excellent' },
            { value: '8', text: '8.0+ Very Good' },
            { value: '7', text: '7.0+ Good' },
            { value: '6', text: '6.0+ Fair' },
            { value: '5', text: '5.0+ Average' }
        ];
        
        ratingSelect.innerHTML = '<option value="all">All Ratings</option>';
        ratings.forEach(rating => {
            const option = document.createElement('option');
            option.value = rating.value;
            option.textContent = rating.text;
            ratingSelect.appendChild(option);
        });
    }

    async function loadRuntimes() {
        const runtimes = [
            { value: '0-90', text: 'Short (< 90 min)' },
            { value: '90-120', text: 'Medium (90-120 min)' },
            { value: '120-180', text: 'Long (120-180 min)' },
            { value: '180-', text: 'Very Long (> 180 min)' }
        ];
        
        runtimeSelect.innerHTML = '<option value="all">Any Length</option>';
        runtimes.forEach(runtime => {
            const option = document.createElement('option');
            option.value = runtime.value;
            option.textContent = runtime.text;
            runtimeSelect.appendChild(option);
        });
    }

    function displayMovies(movies) {
        movieGrid.innerHTML = ''; // Clear existing content
        if (movies && movies.length > 0) {
            movies.forEach(movie => {
                createMovieCard(movie);
            });
        } else {
            movieGrid.innerHTML = '<p>No content found.</p>';
        }
    }

    function appendMovies(movies) {
        if (movies && movies.length > 0) {
            movies.forEach(movie => {
                createMovieCard(movie);
            });
        }
    }

    function createMovieCard(movie) {
        const movieCard = document.createElement('div');
        movieCard.classList.add('movie-card');

        const posterPath = movie.poster_path ? `https://image.tmdb.org/t/p/w200${movie.poster_path}` : 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjMwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjMzMzIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIxNCIgZmlsbD0iI2ZmZiIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPk5vIEltYWdlPC90ZXh0Pjwvc3ZnPg==';
        const title = movie.title || movie.name || 'Untitled';
        const releaseDate = movie.release_date || movie.first_air_date ? new Date(movie.release_date || movie.first_air_date).getFullYear() : 'N/A';
        const voteAverage = movie.vote_average ? movie.vote_average.toFixed(1) : 'N/A';
        const genres = movie.genre_ids ? movie.genre_ids.map(id => getGenreName(id)).join(', ') : 'N/A';
        const isInWatchlist = watchlist.some(item => item.id === movie.id);

        movieCard.innerHTML = `
            <img src="${posterPath}" alt="${title}">
            <div class="movie-info">
                <h3>${title}</h3>
                <p>${releaseDate} &bull; ${genres}</p>
                <div class="movie-actions">
                    <div class="rating"><i class="fas fa-star"></i> ${voteAverage}</div>
                    <button class="watchlist-btn ${isInWatchlist ? 'in-watchlist' : ''}" 
                            onclick="toggleWatchlist(event, ${movie.id}, '${title.replace(/'/g, "\\'")}', '${posterPath}', '${currentContentType}')">
                        ${isInWatchlist ? '✓ In Watchlist' : '+ Watchlist'}
                    </button>
                </div>
            </div>
        `;
        movieCard.dataset.movieId = movie.id;
        movieCard.addEventListener('click', (e) => {
            if (!e.target.classList.contains('watchlist-btn')) {
                fetchMovieDetails(movie.id, currentContentType);
            }
        });
        movieGrid.appendChild(movieCard);
    }

    function showWatchlist() {
        if (watchlist.length === 0) {
            movieGrid.innerHTML = '<p>Your watchlist is empty. Add some movies or TV shows!</p>';
            trendingSection.textContent = 'Your Watchlist';
            return;
        }

        trendingSection.textContent = 'Your Watchlist';
        movieGrid.innerHTML = '';
        
        watchlist.forEach(item => {
            const movieCard = document.createElement('div');
            movieCard.classList.add('movie-card');

            movieCard.innerHTML = `
                <img src="${item.poster}" alt="${item.title}">
                <div class="movie-info">
                    <h3>${item.title}</h3>
                    <p>${item.type === 'movie' ? 'Movie' : 'TV Series'}</p>
                    <div class="movie-actions">
                        <button class="watchlist-btn in-watchlist" 
                                onclick="toggleWatchlist(event, ${item.id}, '${item.title.replace(/'/g, "\\'")}', '${item.poster}', '${item.type}')">
                            ✓ Remove
                        </button>
                    </div>
                </div>
            `;
            movieCard.dataset.movieId = item.id;
            movieCard.addEventListener('click', (e) => {
                if (!e.target.classList.contains('watchlist-btn')) {
                    fetchMovieDetails(item.id, item.type);
                }
            });
            movieGrid.appendChild(movieCard);
        });
    }

    function toggleWatchlist(event, id, title, poster, type) {
        event.stopPropagation();
        
        const existingIndex = watchlist.findIndex(item => item.id === id);
        
        if (existingIndex > -1) {
            // Remove from watchlist
            watchlist.splice(existingIndex, 1);
            event.target.textContent = '+ Watchlist';
            event.target.classList.remove('in-watchlist');
        } else {
            // Add to watchlist
            watchlist.push({ id, title, poster, type });
            event.target.textContent = '✓ In Watchlist';
            event.target.classList.add('in-watchlist');
        }
        
        localStorage.setItem('watchlist', JSON.stringify(watchlist));
        updateWatchlistCount();
    }

    function updateWatchlistCount() {
        const countElement = document.querySelector('.watchlist-count');
        if (countElement) {
            countElement.textContent = `(${watchlist.length})`;
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

    async function fetchMovieDetails(contentId, contentType = 'movie') {
        try {
            let url = `/api/movie/${contentId}`;
            if (contentType === 'tv') {
                url += '?type=tv';
            }
            
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const contentDetails = await response.json();
            displayMovieDetailsModal(contentDetails, contentType);
        } catch (error) {
            console.error('Error fetching content details:', error);
            alert('Failed to load content details. Please try again later.');
        }
    }

    function displayMovieDetailsModal(contentData, contentType = 'movie') {
        console.log('Content Details:', contentData);
        const modal = document.getElementById('movieDetailsModal');
        const modalBody = document.getElementById('modal-body');

        // Extract data from both TMDB and OMDB sources
        const tmdbData = contentData.TMDBData || {};
        const omdbData = contentData.OMDBData || {};

        // Use OMDB data as primary source for detailed info, fallback to TMDB
        const title = omdbData.Title || tmdbData.title || tmdbData.name || 'Unknown Title';
        const year = omdbData.Year || (tmdbData.release_date ? new Date(tmdbData.release_date).getFullYear() : 
                     tmdbData.first_air_date ? new Date(tmdbData.first_air_date).getFullYear() : 'N/A');
        const rated = omdbData.Rated || 'N/A';
        const released = omdbData.Released || tmdbData.release_date || tmdbData.first_air_date || 'N/A';
        
        let runtime = omdbData.Runtime || 'N/A';
        if (contentType === 'tv' && tmdbData.episode_run_time && tmdbData.episode_run_time.length > 0) {
            runtime = `${tmdbData.episode_run_time[0]} min per episode`;
        } else if (tmdbData.runtime) {
            runtime = `${tmdbData.runtime} min`;
        }
        
        const genre = omdbData.Genre || (tmdbData.genres ? tmdbData.genres.map(g => g.name).join(', ') : 'N/A');
        const director = omdbData.Director || (tmdbData.created_by ? tmdbData.created_by.map(c => c.name).join(', ') : 'N/A');
        const writer = omdbData.Writer || 'N/A';
        const actors = omdbData.Actors || 'N/A';
        const plot = omdbData.Plot || tmdbData.overview || 'No plot available';
        const language = omdbData.Language || (tmdbData.original_language || 'N/A');
        const country = omdbData.Country || (tmdbData.production_countries ? tmdbData.production_countries.map(c => c.name).join(', ') : 
                       tmdbData.origin_country ? tmdbData.origin_country.join(', ') : 'N/A');
        const awards = omdbData.Awards || 'N/A';
        const imdbRating = omdbData.imdbRating || (tmdbData.vote_average ? tmdbData.vote_average.toFixed(1) : 'N/A');
        const boxOffice = omdbData.BoxOffice || 'N/A';
        const posterPath = tmdbData.poster_path ? `https://image.tmdb.org/t/p/w300${tmdbData.poster_path}` : 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzAwIiBoZWlnaHQ9IjQ1MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjMzMzIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIxOCIgZmlsbD0iI2ZmZiIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPk5vIEltYWdlPC90ZXh0Pjwvc3ZnPg==';

        // Additional TV show specific info
        let additionalInfo = '';
        if (contentType === 'tv' && tmdbData.number_of_seasons) {
            additionalInfo = `
                <p><strong>Seasons:</strong> ${tmdbData.number_of_seasons}</p>
                <p><strong>Episodes:</strong> ${tmdbData.number_of_episodes || 'N/A'}</p>
                <p><strong>Status:</strong> ${tmdbData.status || 'N/A'}</p>
            `;
        }

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
                    ${additionalInfo}
                    <p><strong>Genre:</strong> ${genre}</p>
                    <p><strong>${contentType === 'tv' ? 'Creator' : 'Director'}:</strong> ${director}</p>
                    <p><strong>Writer:</strong> ${writer}</p>
                    <p><strong>Actors:</strong> ${actors}</p>
                    <p><strong>Plot:</strong> ${plot}</p>
                    <p><strong>Language:</strong> ${language}</p>
                    <p><strong>Country:</strong> ${country}</p>
                    <p><strong>Awards:</strong> ${awards}</p>
                    <p><strong>IMDB Rating:</strong> ${imdbRating}</p>
                    ${contentType === 'movie' ? `<p><strong>Box Office:</strong> ${boxOffice}</p>` : ''}
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

    // Make functions globally available
    window.toggleWatchlist = toggleWatchlist;
});