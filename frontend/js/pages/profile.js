import { loadNavbar } from '../components/navbar.js';
import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  loadNavbar();

  const authState = JSON.parse(localStorage.getItem('authState'));
  if (!authState || !authState.isLoggedIn) {
    alert('Please log in to access your profile.');
    window.location.href = 'login.html';
    return;
  }

  const notificationToggle = document.getElementById('notificationToggle');
  const transactionTableBody = document.getElementById('transactionTable').querySelector('tbody');

  // Fetch and populate the user's notification preferences
  fetch(`${config.apiBaseUrl}/users/${authState.user_id}`, {
    method: 'GET',
    headers: {
      'token': authState.token,
    },
  })
    .then(response => {
      if (!response.ok) throw new Error('Failed to fetch user data.');
      return response.json();
    })
    .then(data => {
      notificationToggle.checked = data.notification_preference === 1;
    })
    .catch(error => {
      console.error('Error fetching notification preferences:', error);
      alert('Failed to load your notification preferences.');
    });

  // Update notification preferences
  notificationToggle.addEventListener('change', () => {
    const notificationsEnabled = notificationToggle.checked;

    fetch(`${config.apiBaseUrl}/users/update/${authState.user_id}`, {
      method: 'PUT',
      headers: {
        'token': authState.token,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ notification_preference: notificationsEnabled ? 1 : 0 }),
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
        alert('Failed to update notification preferences. Reverting changes.');
        notificationToggle.checked = !notificationsEnabled; // Revert on failure
      });
  });

  // Fetch and populate the user's transaction history
  fetch(`${config.apiBaseUrl}/receipts`, {
    method: 'GET',
    headers: {
      'token': authState.token,
    },
  })
    .then(response => {
      if (!response.ok) throw new Error('Failed to fetch transaction history.');
      return response.json();
    })
    .then(data => {
      // Filter transactions for the current user
      const userReceipts = data.filter(receipt => receipt.user_id === authState.user_id);

      if (userReceipts.length === 0) {
        transactionTableBody.innerHTML = '<tr><td colspan="4">No transactions found.</td></tr>';
        return;
      }

      userReceipts.forEach(transaction => {
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
    .catch(error => {
      console.error('Error fetching transaction history:', error);
      alert('Failed to load your transaction history.');
    });
});
