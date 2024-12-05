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
    fetch(`${config.apiBaseUrl}${config.endpoints.reports.revenue}?period=daily`)
      .then(response => response.json())
      .then(data => {
        renderReport(data);
      })
      .catch(error => console.error('Error generating daily report:', error));
  });

  // Generate Monthly Report
  document.getElementById('monthlyReportButton').addEventListener('click', () => {
    fetch(`${config.apiBaseUrl}${config.endpoints.reports.revenue}?period=monthly`)
      .then(response => response.json())
      .then(data => {
        renderReport(data);
      })
      .catch(error => console.error('Error generating monthly report:', error));
  });

  // Render Report
  function renderReport(data) {
    reportOutput.innerHTML = `
      <h2>Report</h2>
      <pre>${JSON.stringify(data, null, 2)}</pre>
    `;
  }
});
