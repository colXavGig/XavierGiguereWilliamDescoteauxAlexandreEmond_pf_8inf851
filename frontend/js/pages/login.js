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
          //extract error message from response
          return response.json().then(data => {
            let error = new Error('Login failed.');
            if (data && data.message) {
              error.message = data.message;
            }
            throw error;
          });
        }
      })
      .then(data => {
        // Save authentication state in localStorage
        const authState = {
          isLoggedIn: true,
          role: data.role, // Assume the backend sends the role
          token: data.token,
        };
        localStorage.setItem('authState', JSON.stringify(authState));
        window.location.href = 'index.html'; // Redirect to the home page
      })
      .catch(error => {
        console.error('Login error:', error);
        errorMessage.textContent = error.message || 'Login failed. Please try again.';
        errorMessage.style.display = 'block';
      });
  });
});
