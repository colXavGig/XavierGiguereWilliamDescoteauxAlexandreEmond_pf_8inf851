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
        <img src="${entity.image_path || 'assets/images/default.jpg'}" alt="${entity.name}">
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

  // Simulate renting an entity
  function rentEntity(entityId) {
    alert(`Entity ${entityId} rented! (Functionality to be implemented.)`);
  }

  // Add filter change listener
  categoryFilter.addEventListener('change', event => {
    fetchEntities(event.target.value);
  });

  // Initial fetch of all entities
  fetchEntities();
});
