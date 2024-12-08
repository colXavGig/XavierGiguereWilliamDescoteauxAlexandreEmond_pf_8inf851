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

  const reportOutput = document.getElementById('reportOutput');

  // Generate Daily Report
  document.getElementById('dailyReportButton').addEventListener('click', () => {
    generateReport('daily');
  });

  // Generate Monthly Report
  document.getElementById('monthlyReportButton').addEventListener('click', () => {
    generateReport('monthly');
  });

  // Generic Report Generation Function
  function generateReport(period) {
    if (!['daily', 'monthly'].includes(period)) {
      alert('Invalid report period specified.');
      return;
    }

    const endpoint = `${config.apiBaseUrl}${config.endpoints.reports.revenue}?period=${period}`;

    fetch(endpoint, {
      headers: {
        'token': authState.token,
      },
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(data => {
            throw new Error(data.message || `Failed to generate ${period} report.`);
          });
        }
        return response.json();
      })
      .then(data => {
        if (!data || typeof data !== 'object') {
          throw new Error('Invalid data format received.');
        }
        renderReport(data, period);
      })
      .catch(error => {
        console.error(`Error generating ${period} report:`, error);
        alert(`Error: ${error.message}`);
      });
  }

  // Render Report
  function renderReport(data, period) {
    reportOutput.innerHTML = `
      <h2>${capitalize(period)} Report</h2>
      <pre>${JSON.stringify(data, null, 2)}</pre>
    `;
  }

  // Utility Function to Capitalize Strings
  function capitalize(str) {
    return str.charAt(0).toUpperCase() + str.slice(1);
  }
});
