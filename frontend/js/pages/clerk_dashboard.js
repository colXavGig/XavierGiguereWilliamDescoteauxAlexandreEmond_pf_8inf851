import { loadNavbar } from '../components/navbar.js';
import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  loadNavbar();

  const authState = JSON.parse(localStorage.getItem('authState'));
  if (!authState || !authState.isLoggedIn || authState.role !== 'clerk') {
    alert('Unauthorized access. Redirecting to login.');
    window.location.href = 'login.html';
    return;
  }

  const rentalLogsTableBody = document.getElementById('rentalLogsTable').querySelector('tbody');

  // Fetch and display rental logs
  fetch(`${config.apiBaseUrl}${config.endpoints.rentalLogs}`)
    .then(response => response.json())
    .then(data => {
      data.forEach(log => {
        const row = document.createElement('tr');
        row.innerHTML = `
          <td>${log.id}</td>
          <td>${log.entity_name}</td>
          <td>${log.user_email}</td>
          <td>${log.rental_date}</td>
          <td>
            <button class="action-button delete" data-id="${log.id}">
              <i class="fas fa-trash"></i> Delete
            </button>
          </td>
        `;
        rentalLogsTableBody.appendChild(row);
      });

      // Add delete event listeners
      document.querySelectorAll('.delete').forEach(button => {
        button.addEventListener('click', event => {
          const logId = event.target.closest('button').dataset.id;
          deleteRentalLog(logId);
        });
      });
    })
    .catch(error => console.error('Error fetching rental logs:', error));

  // Add Rental Log
  document.getElementById('addRentalButton').addEventListener('click', () => {
    alert('Rental Log functionality to be implemented.');
  });

  // Submit Receipts
  document.getElementById('submitReceiptButton').addEventListener('click', () => {
    alert('Submit Receipt functionality to be implemented.');
  });

  // Generate Reports
  document.getElementById('dailyReportButton').addEventListener('click', () => {
    alert('Daily Report functionality to be implemented.');
  });

  document.getElementById('monthlyReportButton').addEventListener('click', () => {
    alert('Monthly Report functionality to be implemented.');
  });

  // Function to delete a rental log
  function deleteRentalLog(logId) {
    fetch(`${config.apiBaseUrl}${config.endpoints.rentalLogs}/${logId}`, {
      method: 'DELETE',
    })
      .then(response => {
        if (response.ok) {
          alert('Rental log deleted successfully.');
          location.reload(); // Refresh the page to update the table
        } else {
          throw new Error('Failed to delete rental log.');
        }
      })
      .catch(error => console.error('Error deleting rental log:', error));
  }
});
