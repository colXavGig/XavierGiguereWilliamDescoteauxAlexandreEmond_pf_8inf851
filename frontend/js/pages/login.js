import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  const loginForm = document.getElementById('loginForm');
  const errorMessage = document.getElementById('errorMessage');

  document.getElementById('loginButton').addEventListener('click', () => {
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value.trim();

    // Clear previous error message
    errorMessage.style.display = 'none';

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
          return response.text().then(text => {
            let error = new Error('Login failed.');
            try {
              const data = JSON.parse(text);
              if (data.message) error.message = data.message;
            } catch {
              error.message = text || 'An unexpected error occurred.';
            }
            throw error;
          });
        }
      })
      .then(data => {
        // Save authentication state in localStorage
        const authState = {
          isLoggedIn: true,
          role: data.user_role, 
          token: data.token,
          user_id: data.user_id, 
        };
        localStorage.setItem('authState', JSON.stringify(authState));
        console.log('Auth state saved:', authState);

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
