import { loadNavbar } from '../components/navbar.js';
import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  loadNavbar();

  const checkoutTableBody = document.getElementById('checkoutTable').querySelector('tbody');
  const totalCostElement = document.getElementById('totalCost');
  const confirmCheckoutButton = document.getElementById('confirmCheckoutButton');
  const clearCartButton = document.getElementById('clearCartButton');

  let cart = JSON.parse(localStorage.getItem('cart')) || [];

  // Validate cart data
  if (!Array.isArray(cart)) {
    localStorage.removeItem('cart');
    alert('Invalid cart data. Your cart has been cleared.');
    cart = [];
  }

  // Render the cart items
  function renderCart() {
    checkoutTableBody.innerHTML = '';
    let totalCost = 0;

    if (cart.length === 0) {
      checkoutTableBody.innerHTML = '<tr><td colspan="4">Your cart is empty.</td></tr>';
      totalCostElement.textContent = 'Total Cost: $0';
      return;
    }

    cart.forEach((item, index) => {
      totalCost += item.price;

      const row = document.createElement('tr');
      row.innerHTML = `
        <td>${item.name}</td>
        <td>${item.category}</td>
        <td>$${item.price.toFixed(2)}</td>
        <td>
          <button class="action-button remove" data-index="${index}">
            <i class="fas fa-times"></i> Remove
          </button>
        </td>
      `;
      checkoutTableBody.appendChild(row);
    });

    totalCostElement.textContent = `Total Cost: $${totalCost.toFixed(2)}`;

    // Add event listeners for remove buttons
    document.querySelectorAll('.remove').forEach(button => {
      button.addEventListener('click', event => {
        const index = event.target.closest('button').dataset.index;
        removeFromCart(index);
      });
    });
  }

  // Remove an item from the cart
  function removeFromCart(index) {
    cart.splice(index, 1);
    localStorage.setItem('cart', JSON.stringify(cart));
    renderCart();
  }

  // Clear the cart
  clearCartButton.addEventListener('click', () => {
    if (confirm('Are you sure you want to clear the cart?')) {
      localStorage.removeItem('cart');
      cart = [];
      renderCart();
    }
  });

  // Confirm checkout
  confirmCheckoutButton.addEventListener('click', () => {
    if (cart.length === 0) {
      alert('Your cart is empty.');
      return;
    }

    const authState = JSON.parse(localStorage.getItem('authState'));
    if (!authState || !authState.isLoggedIn || !authState.token) {
      alert('Your session has expired. Please log in again.');
      window.location.href = 'login.html';
      return;
    }

    const userId = authState.user_id;
    const totalAmount = cart.reduce((sum, item) => sum + item.price, 0);

    const receipt = {
      user_id: userId,
      total_amount: totalAmount,
      status: 'pending',
    };

    // Disable button to prevent duplicate submissions
    confirmCheckoutButton.disabled = true;

    fetch(`${config.apiBaseUrl}${config.endpoints.receipts}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authState.token}`,
      },
      body: JSON.stringify(receipt),
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(data => {
            throw new Error(data.message || 'Failed to complete checkout.');
          });
        }
        return response.json();
      })
      .then(() => {
        alert('Checkout successful!');
        localStorage.removeItem('cart');
        cart = [];
        window.location.href = 'index.html';
      })
      .catch(error => {
        console.error('Error during checkout:', error);
        alert(`Error: ${error.message}`);
      })
      .finally(() => {
        confirmCheckoutButton.disabled = false; // Re-enable the button
      });
  });

  // Initial render
  renderCart();
});
