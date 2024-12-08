import { loadNavbar } from '../components/navbar.js';
import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  loadNavbar();

  const checkoutTableBody = document.getElementById('checkoutTable').querySelector('tbody');
  const totalCostElement = document.getElementById('totalCost');
  const confirmCheckoutButton = document.getElementById('confirmCheckoutButton');
  const clearCartButton = document.getElementById('clearCartButton');

  // Simulating a cart stored in localStorage
  const cart = JSON.parse(localStorage.getItem('cart')) || [];

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
        <td>$${item.price}</td>
        <td>
          <button class="action-button remove" data-index="${index}">
            <i class="fas fa-times"></i> Remove
          </button>
        </td>
      `;
      checkoutTableBody.appendChild(row);
    });

    totalCostElement.textContent = `Total Cost: $${totalCost}`;

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
      renderCart();
    }
  });

  // Confirm checkout
  confirmCheckoutButton.addEventListener('click', () => {
    if (cart.length === 0) {
      alert('Your cart is empty.');
      return;
    }

    // Get user ID from authState
    const authState = JSON.parse(localStorage.getItem('authState'));
    if (!authState || !authState.isLoggedIn) {
      alert('You must be logged in to complete the checkout.');
      window.location.href = 'login.html';
      return;
    }

    const userId = authState.user_id; // Assume `user_id` is stored in authState
    const totalAmount = cart.reduce((sum, item) => sum + item.price, 0); // Calculate total cost

    // Prepare receipt data
    const receipt = {
      user_id: userId,
      total_amount: totalAmount,
      status: 'pending', // Default status
    };

    // Send POST request to submit the receipt
    fetch(`${config.apiBaseUrl}${config.endpoints.receipts}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(receipt),
    })
      .then(response => {
        if (response.ok) {
          alert('Checkout successful!');
          localStorage.removeItem('cart'); // Clear the cart after successful checkout
          window.location.href = 'index.html'; // Redirect to home page
        } else {
          throw new Error('Failed to complete checkout.');
        }
      })
      .catch(error => {
        console.error('Error during checkout:', error);
        alert('An error occurred during checkout. Please try again.');
      });
  });

  // Initial render
  renderCart();
});
