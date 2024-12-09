export function loadNavbar() {
  const authState = JSON.parse(localStorage.getItem('authState')) || { isLoggedIn: false, role: null };
  const navbar = document.getElementById('navbar');

  navbar.innerHTML = `
    <nav class="navbar">
      <!-- Common Links -->
      <a href="index.html"><i class="fas fa-home"></i> Home</a>
      ${authState.isLoggedIn ? `
        <a href="profile.html"><i class="fas fa-user"></i> Profile</a>
        ${authState.role === 'clerk' || authState.role === 'admin' ? `
          <a href="clerk_dashboard.html"><i class="fas fa-tools"></i> Dashboard</a>
        ` : ''}
        ${authState.role === 'admin' ? `
          <a href="admin_user_management.html"><i class="fas fa-users"></i> User Management</a>
          <a href="admin_receipts.html"><i class="fas fa-file-invoice-dollar"></i> Receipts</a>
          <a href="admin_reports.html"><i class="fas fa-chart-bar"></i> Reports</a>
        ` : ''}
        <a id="logout" href="#"><i class="fas fa-sign-out-alt"></i> Logout</a>
      ` : `
        <a href="login.html"><i class="fas fa-sign-in-alt"></i> Login</a>
      `}
    </nav>
  `;

  // Logout logic
  const logoutButton = document.getElementById('logout');
  if (logoutButton) {
    logoutButton.addEventListener('click', () => {
      localStorage.removeItem('authState');
      window.location.reload();
    });
  }
}
