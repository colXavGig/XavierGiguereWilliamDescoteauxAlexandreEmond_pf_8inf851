document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');
    const errorMessage = document.getElementById('errorMessage');
  
    document.getElementById('registerButton').addEventListener('click', () => {
      const email = document.getElementById('email').value;
      const password = document.getElementById('password').value;
      const confirmPassword = document.getElementById('confirmPassword').value;
  
      if (password !== confirmPassword) {
        errorMessage.textContent = 'Passwords do not match. Please try again.';
        errorMessage.style.display = 'block';
        return;
      }
  
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
            errorMessage.textContent = 'Registration failed. Please try again.';
            errorMessage.style.display = 'block';
          }
        })
        .catch(error => {
          console.error('Error during registration:', error);
          errorMessage.textContent = 'An error occurred. Please try again.';
          errorMessage.style.display = 'block';
        });
    });
  });
  