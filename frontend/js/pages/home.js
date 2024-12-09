import { loadNavbar } from '../components/navbar.js';
import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  loadNavbar();

  const entityGrid = document.getElementById('entityGrid');
  const categoryFilter = document.getElementById('categoryFilter');
  const authState = JSON.parse(localStorage.getItem('authState'));

  // Check if the user is logged in
  if (!authState || !authState.isLoggedIn || !authState.token) {
    alert('Your session has expired. Please log in again.');
    window.location.href = 'login.html';
    return;
  }

  const token = authState.token; // Use token from authState

  // Fetch and display all rentable entities
  function fetchEntities(filterCategory = 'all') {
    const url =
      filterCategory === 'all'
        ? `${config.apiBaseUrl}${config.endpoints.rentableEntities}`
        : `${config.apiBaseUrl}${config.endpoints.rentableEntities}?category=${filterCategory}`;

    entityGrid.innerHTML = '<p>Loading entities...</p>';

    fetch(url, {
      headers: {
        'token': token,
      },
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Failed to fetch entities.');
        }
        return response.json();
      })
      .then(data => {
        renderEntities(data);
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
    document.querySelectorAll('.rent').forEach(button => {
      button.addEventListener('click', event => {
        const entityId = event.target.closest('button').dataset.id;
        rentEntity(entityId);
      });
    });
  }

  // Rent an entity
  function rentEntity(entityId) {
    const rentalData = {
      entity_id: entityId,
      user_id: authState.user_id,
      rental_date: new Date().toISOString().split('T')[0],
    };

    fetch(`${config.apiBaseUrl}${config.endpoints.rentalLogs}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'token': token,
      },
      body: JSON.stringify(rentalData),
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(error => {
            throw new Error(error.message || 'Failed to rent entity.');
          });
        }
        alert('Rental successful!');
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
