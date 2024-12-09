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

  const usersTableBody = document.getElementById('usersTable').querySelector('tbody');
  const userForm = document.getElementById('userForm');
  const userFormModal = new bootstrap.Modal(document.getElementById('userFormModal'));
  const saveUserButton = document.getElementById('saveUserButton');

  let editingUserId = null; // To track whether we are adding or editing a user

  // Fetch and display all users
  function fetchUsers() {
    fetch(`${config.apiBaseUrl}${config.endpoints.users}`, {
      headers: {
        'token': authState.token,
      },
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(data => {
            throw new Error(data.message || 'Failed to fetch users.');
          });
        }
        return response.json();
      })
      .then(users => {
        renderUsersTable(users);
      })
      .catch(error => {
        console.error('Error fetching users:', error);
        alert(`Error: ${error.message}`);
      });
  }

  // Render the users table
  function renderUsersTable(users) {
    usersTableBody.innerHTML = '';
    if (users.length === 0) {
      usersTableBody.innerHTML = '<tr><td colspan="5">No users found.</td></tr>';
      return;
    }

    users.forEach(user => {
      const row = document.createElement('tr');
      row.innerHTML = `
        <td>${user.id}</td>
        <td>${user.email}</td>
        <td>${user.role}</td>
        <td>
          <button class="btn btn-primary btn-sm edit" data-id="${user.id}">
            <i class="fas fa-edit"></i> Edit
          </button>
          <button class="btn btn-danger btn-sm delete" data-id="${user.id}">
            <i class="fas fa-trash"></i> Delete
          </button>
        </td>
      `;
      usersTableBody.appendChild(row);
    });

    // Add event listeners for edit and delete buttons
    document.querySelectorAll('.edit').forEach(button => {
      button.addEventListener('click', event => {
        const userId = event.target.closest('button').dataset.id;
        openEditUserModal(userId);
      });
    });

    document.querySelectorAll('.delete').forEach(button => {
      button.addEventListener('click', event => {
        const userId = event.target.closest('button').dataset.id;
        deleteUser(userId);
      });
    });
  }

  // Open the user form modal for editing
  function openEditUserModal(userId) {
    editingUserId = userId;

    fetch(`${config.apiBaseUrl}${config.endpoints.users}/${userId}`, {
      headers: {
        'token': authState.token,
      },
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(data => {
            throw new Error(data.message || 'Failed to fetch user details.');
          });
        }
        return response.json();
      })
      .then(user => {
        document.getElementById('email').value = user.email;
        document.getElementById('role').value = user.role;
        userFormModal.show();
      })
      .catch(error => {
        console.error('Error fetching user details:', error);
        alert(`Error: ${error.message}`);
      });
  }

  // Save user (create or update)
  saveUserButton.addEventListener('click', () => {
    const email = document.getElementById('email').value.trim();
    const role = document.getElementById('role').value;

    if (!email || !role) {
      alert('All fields are required.');
      return;
    }

    const userData = { email, role };

    const endpoint = editingUserId
      ? `${config.apiBaseUrl}${config.endpoints.users}/${editingUserId}`
      : `${config.apiBaseUrl}${config.endpoints.users}/create`;
    const method = editingUserId ? 'PUT' : 'POST';

    fetch(endpoint, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'token': authState.token,
      },
      body: JSON.stringify(userData),
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(data => {
            throw new Error(data.message || 'Failed to save user.');
          });
        }
        userFormModal.hide();
        fetchUsers();
        alert(`User ${editingUserId ? 'updated' : 'created'} successfully.`);
        editingUserId = null; // Reset editing ID
      })
      .catch(error => {
        console.error('Error saving user:', error);
        alert(`Error: ${error.message}`);
      });
  });

  // Delete user
  function deleteUser(userId) {
    if (!confirm('Are you sure you want to delete this user?')) {
      return;
    }

    fetch(`${config.apiBaseUrl}${config.endpoints.users}/${userId}`, {
      method: 'DELETE',
      headers: {
        'token': authState.token,
      },
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(data => {
            throw new Error(data.message || 'Failed to delete user.');
          });
        }
        fetchUsers();
        alert('User deleted successfully.');
      })
      .catch(error => {
        console.error('Error deleting user:', error);
        alert(`Error: ${error.message}`);
      });
  }

  // Fetch users on page load
  fetchUsers();
});
