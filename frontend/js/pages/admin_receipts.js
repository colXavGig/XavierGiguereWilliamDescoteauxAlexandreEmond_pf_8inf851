import { loadNavbar } from '../components/navbar.js';
import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  loadNavbar();

  const authState = JSON.parse(localStorage.getItem('authState'));
  if (!authState || !authState.isLoggedIn || authState.role !== 'admin') {
    alert('Unauthorized access. Redirecting to login.');
    window.location.href = 'login.html';
    return;
  }

  const receiptsTableBody = document.getElementById('receiptsTable').querySelector('tbody');

  // Fetch and display receipts
  fetch(`${config.apiBaseUrl}${config.endpoints.receipts}`)
    .then(response => response.json())
    .then(data => {
      data.forEach(receipt => {
        const row = document.createElement('tr');
        row.innerHTML = `
          <td>${receipt.id}</td>
          <td>${receipt.user_email}</td>
          <td>$${receipt.total_amount}</td>
          <td>${receipt.status}</td>
          <td>
            ${receipt.status === 'pending' ? `
              <button class="action-button approve" data-id="${receipt.id}">
                <i class="fas fa-check"></i> Approve
              </button>
              <button class="action-button reject" data-id="${receipt.id}">
                <i class="fas fa-times"></i> Reject
              </button>
            ` : ''}
          </td>
        `;
        receiptsTableBody.appendChild(row);
      });

      // Add event listeners for Approve and Reject buttons
      document.querySelectorAll('.approve').forEach(button => {
        button.addEventListener('click', event => {
          const receiptId = event.target.closest('button').dataset.id;
          updateReceiptStatus(receiptId, 'approved');
        });
      });

      document.querySelectorAll('.reject').forEach(button => {
        button.addEventListener('click', event => {
          const receiptId = event.target.closest('button').dataset.id;
          updateReceiptStatus(receiptId, 'rejected');
        });
      });
    })
    .catch(error => console.error('Error fetching receipts:', error));

  // Function to update receipt status
  function updateReceiptStatus(receiptId, status) {
    fetch(`${config.apiBaseUrl}${config.endpoints.receipts}/${receiptId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status }),
    })
      .then(response => {
        if (response.ok) {
          alert(`Receipt ${status} successfully.`);
          location.reload(); // Refresh the page to update the table
        } else {
          throw new Error(`Failed to ${status} receipt.`);
        }
      })
      .catch(error => console.error(`Error updating receipt:`, error));
  }
});
