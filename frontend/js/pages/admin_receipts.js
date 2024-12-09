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
  fetch(`${config.apiBaseUrl}${config.endpoints.receipts}`, {
    headers: {
      'token': authState.token,
    },
  })
    .then(response => {
      if (!response.ok) {
        return response.json().then(data => {
          throw new Error(data.message || 'Failed to fetch receipts.');
        });
      }
      return response.json();
    })
    .then(data => {
      if (!Array.isArray(data)) {
        throw new Error('Invalid data format received.');
      }

      if (data.length === 0) {
        receiptsTableBody.innerHTML = '<tr><td colspan="5">No receipts found.</td></tr>';
        return;
      }

      data.forEach(receipt => {
        const row = document.createElement('tr');
        row.innerHTML = `
          <td>${receipt.id}</td>
          <td>${receipt.user_email || 'Unknown'}</td>
          <td>$${receipt.total_amount.toFixed(2)}</td>
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
    .catch(error => {
      console.error('Error fetching receipts:', error);
      alert(`Error: ${error.message}`);
    });

  // Function to update receipt status
  function updateReceiptStatus(receiptId, status) {
    fetch(`${config.apiBaseUrl}${config.endpoints.receipts}/update/${receiptId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'token': authState.token,
      },
      body: JSON.stringify({ status }),
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(data => {
            throw new Error(data.message || `Failed to ${status} receipt.`);
          });
        }
        // Update status dynamically in the table
        const row = document.querySelector(`button[data-id="${receiptId}"]`).closest('tr');
        row.querySelector('td:nth-child(4)').textContent = status;
        row.querySelector('td:nth-child(5)').innerHTML = ''; // Clear action buttons
        alert(`Receipt ${status} successfully.`);
      })
      .catch(error => {
        console.error(`Error updating receipt:`, error);
        alert(`Error: ${error.message}`);
      });
  }
});
