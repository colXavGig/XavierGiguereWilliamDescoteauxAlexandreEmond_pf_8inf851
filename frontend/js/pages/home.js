import { loadNavbar } from '../components/navbar.js';
import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  loadNavbar();

  const entityGrid = document.getElementById('entityGrid');
  const categoryFilter = document.getElementById('categoryFilter');
  const authState = JSON.parse(localStorage.getItem('authState'));

  // Log authState for debugging
  console.log('Auth state:', authState);

  // Check if the user is logged in
  if (!authState || !authState.isLoggedIn || !authState.token) {
    alert('Your session has expired. Please log in again.');
    window.location.href = 'login.html';
    return;
  }

  const token = authState.token; // Use token from authState
  let allEntities = []; // Store all entities fetched from the API

  /**
   * Fetch all entities initially and store them in memory.
   */
  function fetchEntities() {
    const url = `${config.apiBaseUrl}${config.endpoints.rentableEntities}`;

    // Show loading indicator
    entityGrid.innerHTML = '<p>Loading entities...</p>';

    fetch(url, {
      headers: {
        'token': token,
      },
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error('Failed to fetch entities.');
        }
        return response.json();
      })
      .then((data) => {
        allEntities = data; // Store fetched entities
        renderEntities(allEntities); // Render all entities initially
      })
      .catch((error) => {
        console.error('Error fetching entities:', error);
        entityGrid.innerHTML = '<p>Error loading entities. Please try again later.</p>';
      });
  }

  /**
   * Render entities in the grid, filtered by category if provided.
   * @param {Array} entities - The list of entities to display.
   * @param {string} filterCategory - The category to filter by. Defaults to 'all'.
   */
  function renderEntities(entities, filterCategory = 'all') {
    entityGrid.innerHTML = '';

    // Apply the filter if it's not 'all'
    const filteredEntities =
      filterCategory.toLowerCase() === 'all'
        ? entities
        : entities.filter((entity) =>
            entity.category.toLowerCase() === filterCategory.toLowerCase()
          );

    // Check if there are no results
    if (filteredEntities.length === 0) {
      entityGrid.innerHTML = '<p>No entities available for this category.</p>';
      return;
    }

    // Render the filtered entities
    filteredEntities.forEach((entity) => {
      const card = document.createElement('div');
      card.className = 'entity-card';
      card.innerHTML = `
        <img src="${entity.image_path || 'images/default.jpg'}" alt="${entity.name}">
        <h3>${entity.name}</h3>
        <p>${entity.description || 'No description available.'}</p>
        <p><strong>Price:</strong> $${entity.price} (${entity.pricing_model})</p>
        <button class="action-button rent" data-id="${entity.id}">
          <i class="fas fa-plus-circle"></i> Rent
        </button>
      `;
      entityGrid.appendChild(card);
    });

    // Add event listeners for "Rent" buttons
    document.querySelectorAll('.rent').forEach((button) => {
      button.addEventListener('click', (event) => {
        const entityId = event.target.closest('button').dataset.id;
        rentEntity(entityId);
      });
    });
  }

  /**
   * Handle the renting of an entity.
   * @param {string} entityId - The ID of the entity to rent.
   */
  function rentEntity(entityId) {
    const rentalData = {
      entity_id: entityId,
      user_id: authState.user_id,
      rental_date: new Date().toISOString().split('T')[0], // Current date
    };

    fetch(`${config.apiBaseUrl}${config.endpoints.rentalLogs}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'token': token,
      },
      body: JSON.stringify(rentalData),
    })
      .then((response) => {
        if (!response.ok) {
          return response.json().then((error) => {
            throw new Error(error.message || 'Failed to rent entity.');
          });
        }
        alert('Rental successful!');
      })
      .catch((error) => {
        console.error('Error renting entity:', error);
        alert('An error occurred while renting the entity. Please try again.');
      });
  }

  /**
   * Handle the filtering logic when the category changes.
   */
  categoryFilter.addEventListener('change', (event) => {
    const selectedCategory = event.target.value;
    renderEntities(allEntities, selectedCategory); // Filter entities in the frontend
  });

  // Fetch all entities on initial load
  fetchEntities();
});
