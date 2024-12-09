import config from '../config.js';

document.addEventListener('DOMContentLoaded', () => {
  const registerForm = document.getElementById('registerForm');
  const errorMessage = document.getElementById('errorMessage');

  document.getElementById('registerButton').addEventListener('click', () => {
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value.trim();
    const confirmPassword = document.getElementById('confirmPassword').value.trim();

    // Clear previous error messages
    errorMessage.style.display = 'none';

    // Validate email format
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      errorMessage.textContent = 'Please enter a valid email address.';
      errorMessage.style.display = 'block';
      return;
    }

    // Validate password length
    if (password.length < 8) {
      errorMessage.textContent = 'Password must be at least 8 characters long.';
      errorMessage.style.display = 'block';
      return;
    }

    // Validate password confirmation
    if (password !== confirmPassword) {
      errorMessage.textContent = 'Passwords do not match. Please try again.';
      errorMessage.style.display = 'block';
      return;
    }

    // Send registration request
    fetch(`${config.apiBaseUrl}${config.endpoints.register}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    })
      .then(response => {
        if (response.ok) {
          alert('Registration successful! You can now log in.');
          window.location.href = 'login.html'; // Redirect to login page
        } else {
          return response.json().then(data => {
            errorMessage.textContent = data.message || 'Registration failed. Please try again.';
            errorMessage.style.display = 'block';
          });
        }
      })
      .catch(error => {
        console.error('Error during registration:', error);
        errorMessage.textContent = 'An error occurred. Please try again.';
        errorMessage.style.display = 'block';
      });
  });
});
