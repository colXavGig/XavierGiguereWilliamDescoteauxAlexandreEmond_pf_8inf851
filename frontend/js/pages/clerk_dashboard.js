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
  fetch(`${config.apiBaseUrl}${config.endpoints.rentalLogs}`, {
    headers: {
      'Authorization': `Bearer ${authState.token}`,
    },
  })
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
    .catch(error => {
      console.error('Error fetching rental logs:', error);
      alert('Failed to load rental logs.');
    });

  // Function to add a rental log
  function addRentalLog(entityId, userId, rentalDate) {
    if (!entityId || !userId || !rentalDate) {
      alert('All fields are required.');
      return;
    }

    fetch(`${config.apiBaseUrl}${config.endpoints.rentalLogs}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authState.token}`,
      },
      body: JSON.stringify({
        entity_id: parseInt(entityId),
        user_id: parseInt(userId),
        rental_date: rentalDate,
      }),
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(data => {
            throw new Error(data.message || 'Failed to add rental log.');
          });
        }
        return response.json();
      })
      .then(log => {
        // Dynamically add the new log to the table
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
        alert('Rental log added successfully.');
      })
      .catch(error => {
        console.error('Error adding rental log:', error);
        alert(`Error: ${error.message}`);
      });
  }

  // Function to delete a rental log
  function deleteRentalLog(logId) {
    fetch(`${config.apiBaseUrl}${config.endpoints.rentalLogs}/${logId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${authState.token}`,
      },
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(data => {
            throw new Error(data.message || 'Failed to delete rental log.');
          });
        }
        // Dynamically remove the log from the table
        document.querySelector(`button[data-id="${logId}"]`).closest('tr').remove();
        alert('Rental log deleted successfully.');
      })
      .catch(error => {
        console.error('Error deleting rental log:', error);
        alert(`Error: ${error.message}`);
      });
  }

  // Function to submit a receipt
function submitReceipt(userId, totalAmount) {
  if (!userId || isNaN(userId)) {
    alert('Invalid User ID. Please enter a valid numeric User ID.');
    return;
  }

  if (!totalAmount || isNaN(totalAmount) || parseFloat(totalAmount) <= 0) {
    alert('Invalid Total Amount. Please enter a positive number.');
    return;
  }

  fetch(`${config.apiBaseUrl}${config.endpoints.receipts}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${authState.token}`,
    },
    body: JSON.stringify({
      user_id: parseInt(userId),
      total_amount: parseFloat(totalAmount),
      status: 'pending', // Assuming status defaults to 'pending'
    }),
  })
    .then(response => {
      if (!response.ok) {
        return response.json().then(data => {
          throw new Error(data.message || 'Failed to submit receipt.');
        });
      }
      return response.json();
    })
    .then(data => {
      alert('Receipt submitted successfully.');
      console.log('Submitted receipt:', data);
      // Optionally update UI if needed
    })
    .catch(error => {
      console.error('Error submitting receipt:', error);
      alert(`Error: ${error.message}`);
    });
}
// Function to generate a report
function generateReport(type, dateOrMonth) {
  if (type === 'daily') {
    // Validate date format (YYYY-MM-DD)
    const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
    if (!dateRegex.test(dateOrMonth)) {
      alert('Invalid date format. Please use YYYY-MM-DD.');
      return;
    }
  } else if (type === 'monthly') {
    // Validate month format (YYYY-MM)
    const monthRegex = /^\d{4}-\d{2}$/;
    if (!monthRegex.test(dateOrMonth)) {
      alert('Invalid month format. Please use YYYY-MM.');
      return;
    }
  } else {
    alert('Invalid report type.');
    return;
  }

  const endpoint =
    type === 'daily'
      ? `${config.apiBaseUrl}${config.endpoints.reports.revenue}?date=${dateOrMonth}`
      : `${config.apiBaseUrl}${config.endpoints.reports.availability}?month=${dateOrMonth}`;

  fetch(endpoint, {
    headers: {
      'Authorization': `Bearer ${authState.token}`,
    },
  })
    .then(response => {
      if (!response.ok) {
        return response.json().then(data => {
          throw new Error(data.message || 'Failed to generate report.');
        });
      }
      return response.json();
    })
    .then(data => {
      alert(`Report generated successfully. Check console for details.`);
      console.log(`Generated ${type} report:`, data);
      // Optionally display report data in the UI
    })
    .catch(error => {
      console.error(`Error generating ${type} report:`, error);
      alert(`Error: ${error.message}`);
    });
}

});
