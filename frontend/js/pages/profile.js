import { loadNavbar } from '../components/navbar.js';
import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  loadNavbar();

  // Fetch the user's profile data
  const authState = JSON.parse(localStorage.getItem('authState'));
  if (!authState || !authState.isLoggedIn) {
    alert('Please log in to access your profile.');
    window.location.href = 'login.html';
    return;
  }

  const notificationToggle = document.getElementById('notificationToggle');
  const transactionTableBody = document.getElementById('transactionTable').querySelector('tbody');

  // Fetch and populate the user's notification preferences
  fetch(`${config.apiBaseUrl}/users/${authState.token}/preferences`)
    .then(response => response.json())
    .then(data => {
      notificationToggle.checked = data.notificationsEnabled;
    })
    .catch(error => console.error('Error fetching notification preferences:', error));

  // Update notification preferences
  notificationToggle.addEventListener('change', () => {
    const notificationsEnabled = notificationToggle.checked;

    fetch(`${config.apiBaseUrl}/users/${authState.token}/preferences`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ notificationsEnabled }),
    })
      .then(response => {
        if (response.ok) {
          alert('Notification preferences updated.');
        } else {
          throw new Error('Failed to update preferences.');
        }
      })
      .catch(error => {
        console.error('Error updating preferences:', error);
        notificationToggle.checked = !notificationsEnabled; // Revert on failure
      });
  });

  // Fetch and populate the user's transaction history
  fetch(`${config.apiBaseUrl}/receipts?user_id=${authState.token}`)
    .then(response => response.json())
    .then(data => {
      data.forEach(transaction => {
        const row = document.createElement('tr');
        row.innerHTML = `
          <td>${transaction.id}</td>
          <td>${transaction.date}</td>
          <td>$${transaction.amount}</td>
          <td>${transaction.status}</td>
        `;
        transactionTableBody.appendChild(row);
      });
    })
    .catch(error => console.error('Error fetching transaction history:', error));
});
