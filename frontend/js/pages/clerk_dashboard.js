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
    const entityId = prompt('Enter Entity ID:');
    const userId = prompt('Enter User ID:');
    const rentalDate = prompt('Enter Rental Date (YYYY-MM-DD):');

    if (entityId && userId && rentalDate) {
      addRentalLog(entityId, userId, rentalDate);
    } else {
      alert('All fields are required to add a rental log.');
    }
  });

  // Submit Receipts
  document.getElementById('submitReceiptButton').addEventListener('click', () => {
    const userId = prompt('Enter User ID:');
    const totalAmount = prompt('Enter Total Amount:');

    if (userId && totalAmount) {
      submitReceipt(userId, totalAmount);
    } else {
      alert('All fields are required to submit a receipt.');
    }
  });

  // Generate Reports
  document.getElementById('dailyReportButton').addEventListener('click', () => {
    const date = prompt('Enter Date for Daily Report (YYYY-MM-DD):');
    if (date) {
      generateReport('daily', date);
    } else {
      alert('Date is required to generate a daily report.');
    }
  });

  document.getElementById('monthlyReportButton').addEventListener('click', () => {
    const month = prompt('Enter Month for Monthly Report (YYYY-MM):');
    if (month) {
      generateReport('monthly', month);
    } else {
      alert('Month is required to generate a monthly report.');
    }
  });

  // Function to add a rental log
  function addRentalLog(entityId, userId, rentalDate) {
    fetch(`${config.apiBaseUrl}${config.endpoints.rentalLogs}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        entity_id: entityId,
        user_id: userId,
        rental_date: rentalDate,
      }),
    })
      .then(response => {
        if (response.ok) {
          alert('Rental log added successfully.');
          location.reload(); // Refresh the page to update the table
        } else {
          throw new Error('Failed to add rental log.');
        }
      })
      .catch(error => console.error('Error adding rental log:', error));
  }

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

  // Function to submit a receipt
  function submitReceipt(userId, totalAmount) {
    fetch(`${config.apiBaseUrl}${config.endpoints.receipts}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        user_id: userId,
        total_amount: parseFloat(totalAmount),
        status: 'pending',
      }),
    })
      .then(response => {
        if (response.ok) {
          alert('Receipt submitted successfully.');
          location.reload();
        } else {
          throw new Error('Failed to submit receipt.');
        }
      })
      .catch(error => console.error('Error submitting receipt:', error));
  }

  // Function to generate a report
  function generateReport(type, dateOrMonth) {
    const endpoint =
      type === 'daily'
        ? `${config.apiBaseUrl}${config.endpoints.reports.daily}?date=${dateOrMonth}`
        : `${config.apiBaseUrl}${config.endpoints.reports.monthly}?month=${dateOrMonth}`;

    fetch(endpoint)
      .then(response => response.json())
      .then(data => {
        alert(`Report generated: ${JSON.stringify(data)}`);
        console.log(`Generated ${type} report:`, data);
      })
      .catch(error => console.error(`Error generating ${type} report:`, error));
  }
});
