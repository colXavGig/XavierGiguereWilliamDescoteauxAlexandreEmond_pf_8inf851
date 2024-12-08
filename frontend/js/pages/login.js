import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  const loginForm = document.getElementById('loginForm');
  const errorMessage = document.getElementById('errorMessage');

  document.getElementById('loginButton').addEventListener('click', () => {
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value.trim();

    if (!email || !password) {
      errorMessage.textContent = 'Please enter both email and password.';
      errorMessage.style.display = 'block';
      return;
    }

    // Send login request to the backend
    fetch(`${config.apiBaseUrl}${config.endpoints.login}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    })
      .then(response => {
        if (response.ok) {
          return response.json();
        } else {
          // Parse error response if possible
          return response.json().then(error => {
            throw new Error(error.message || 'Invalid email or password.');
          });
        }
      })
      .then(data => {
        // Validate the response structure
        if (!data.token || !data.role || !data.user_id) {
          throw new Error('Invalid response from server. Please try again.');
        }

        // Save authentication state in localStorage
        const authState = {
          isLoggedIn: true,
          user_id: data.user_id, // Backend should return user_id
          role: data.role,       // Backend should return role
          token: data.token,     // Backend should return a JWT token
        };
        localStorage.setItem('authState', JSON.stringify(authState));

        // Redirect to the home page
        window.location.href = 'index.html';
      })
      .catch(error => {
        console.error('Login error:', error);
        errorMessage.textContent = error.message || 'Login failed. Please try again.';
        errorMessage.style.display = 'block';
      });
  });
});
