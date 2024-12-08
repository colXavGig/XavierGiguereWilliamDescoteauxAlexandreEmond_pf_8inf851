import { loadNavbar } from '../components/navbar.js';
import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  loadNavbar();

  const entityGrid = document.getElementById('entityGrid');
  const categoryFilter = document.getElementById('categoryFilter');

  // Fetch and display all rentable entities
  function fetchEntities(filterCategory = 'all') {
    const url = `${config.apiBaseUrl}${config.endpoints.rentableEntities}`;
    entityGrid.innerHTML = '<p>Loading entities...</p>';

    fetch(url)
      .then(response => response.json())
      .then(data => {
        // Filter by category if a specific category is selected
        const filteredData = filterCategory === 'all'
          ? data
          : data.filter(entity => entity.category === filterCategory);

        // Render the entities
        renderEntities(filteredData);
      })
      .catch(error => {
        console.error('Error fetching entities:', error);
        entityGrid.innerHTML = '<p>Error loading entities. Please try again later.</p>';
      });
  }

  // Render entities in the grid
  function renderEntities(entities) {
    entityGrid.innerHTML = '';

    if (entities.length === 0) {
      entityGrid.innerHTML = '<p>No entities available for this category.</p>';
      return;
    }

    entities.forEach(entity => {
      const card = document.createElement('div');
      card.className = 'entity-card';
      card.innerHTML = `
        <img src="${entity.image_path || 'images/magasin1.jpg'}" alt="${entity.name}">
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
    document.querySelectorAll('.rent').forEach(button => {
      button.addEventListener('click', event => {
        const entityId = event.target.closest('button').dataset.id;
        rentEntity(entityId);
      });
    });
  }

  // Rent an entity
  function rentEntity(entityId) {
    // Check if user is logged in
    const authState = JSON.parse(localStorage.getItem('authState'));
    if (!authState || !authState.isLoggedIn) {
      alert('You must be logged in to rent an entity.');
      window.location.href = 'login.html';
      return;
    }

    // Prepare rental data
    const rentalData = {
      entity_id: entityId,
      user_id: authState.user_id, // Assuming user_id is in authState
      rental_date: new Date().toISOString().split('T')[0], // Current date in YYYY-MM-DD format
    };

    // Send POST request to register the rental
    fetch(`${config.apiBaseUrl}${config.endpoints.rentalLogs}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(rentalData),
    })
      .then(response => {
        if (response.ok) {
          alert('Rental successful!');
          // Optional: Refresh the grid or display a success message
        } else {
          return response.json().then(error => {
            throw new Error(error.message || 'Failed to rent entity.');
          });
        }
      })
      .catch(error => {
        console.error('Error renting entity:', error);
        alert('An error occurred while renting the entity. Please try again.');
      });
  }

  // Add filter change listener
  categoryFilter.addEventListener('change', event => {
    fetchEntities(event.target.value);
  });

  // Initial fetch of all entities
  fetchEntities();
});
