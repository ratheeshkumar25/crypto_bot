// Login page JavaScript
document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginFormElement');
    const loginBtn = document.querySelector('.register-btn');

    // Check if already logged in
    const token = localStorage.getItem('token');
    if (token) {
        // Redirect to dashboard if already logged in
        window.location.href = '/dashboard';
        return;
    }

    // Login form submission
    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const username = document.getElementById('loginUsername').value;
        const password = document.getElementById('loginPassword').value;

        if (!username || !password) {
            showError('Please fill in all fields');
            return;
        }

        loginBtn.classList.add('loading');
        loginBtn.textContent = 'Signing In...';

        try {
                                    const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password }),
            });

            const data = await response.json();

            if (response.ok) {
                // Store token and user data
                localStorage.setItem('token', data.token);
                localStorage.setItem('user', JSON.stringify(data.user));

                // Show success message
                showSuccess('Login successful! Redirecting to dashboard...');

                // Redirect after a short delay
                setTimeout(() => {
window.location.href = '/dashboard';
                }, 500);
            } else {
                showError(data.error || 'Login failed');
            }
        } catch (error) {
            showError('Network error: ' + error.message);
        } finally {
            loginBtn.classList.remove('loading');
            loginBtn.textContent = 'Sign In';
        }
    });

    function showError(message) {
        // Remove existing messages
        document.querySelectorAll('.error-message, .success-message').forEach(el => el.remove());

        const errorDiv = document.createElement('div');
        errorDiv.className = 'error-message';
        errorDiv.textContent = message;
        errorDiv.style.cssText = `
            background: #f8d7da;
            color: #721c24;
            padding: 10px;
            border-radius: 5px;
            margin-top: 20px;
            text-align: center;
            border: 1px solid #f5c6cb;
        `;

        loginForm.appendChild(errorDiv);
    }

    function showSuccess(message) {
        // Remove existing messages
        document.querySelectorAll('.error-message, .success-message').forEach(el => el.remove());

        const successDiv = document.createElement('div');
        successDiv.className = 'success-message';
        successDiv.textContent = message;
        successDiv.style.cssText = `
            background: #d4edda;
            color: #155724;
            padding: 10px;
            border-radius: 5px;
            margin-top: 20px;
            text-align: center;
            border: 1px solid #c3e6cb;
        `;

        loginForm.appendChild(successDiv);
    }
});
