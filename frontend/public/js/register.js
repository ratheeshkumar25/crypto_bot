// Register page JavaScript
document.addEventListener('DOMContentLoaded', function() {
    const registerForm = document.getElementById('registerFormElement');
    const registerBtn = document.querySelector('.register-btn');

    // Exchange integration toggles
    const enableBinance = document.getElementById('enableBinance');
    const binanceKeys = document.getElementById('binanceKeys');
    const enableSolana = document.getElementById('enableSolana');
    const solanaKeys = document.getElementById('solanaKeys');

    // Toggle exchange key fields
    enableBinance.addEventListener('change', function() {
        binanceKeys.style.display = this.checked ? 'block' : 'none';
    });

    enableSolana.addEventListener('change', function() {
        solanaKeys.style.display = this.checked ? 'block' : 'none';
    });

    // Form validation
    function validateForm() {
        const username = document.getElementById('registerUsername').value;
        const email = document.getElementById('registerEmail').value;
        const password = document.getElementById('registerPassword').value;
        const confirmPassword = document.getElementById('confirmPassword').value;

        // Clear previous errors
        document.querySelectorAll('.error-message').forEach(el => el.remove());
        document.querySelectorAll('.form-group').forEach(el => el.classList.remove('error'));

        let isValid = true;

        // Username validation
        if (username.length < 3) {
            showError('registerUsername', 'Username must be at least 3 characters');
            isValid = false;
        }

        // Email validation
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email)) {
            showError('registerEmail', 'Please enter a valid email address');
            isValid = false;
        }

        // Password validation
        if (password.length < 6) {
            showError('registerPassword', 'Password must be at least 6 characters');
            isValid = false;
        }

        // Confirm password
        if (password !== confirmPassword) {
            showError('confirmPassword', 'Passwords do not match');
            isValid = false;
        }

        return isValid;
    }

    function showError(inputId, message) {
        const input = document.getElementById(inputId);
        const formGroup = input.closest('.form-group');
        formGroup.classList.add('error');

        const errorDiv = document.createElement('span');
        errorDiv.className = 'error-message';
        errorDiv.textContent = message;
        formGroup.appendChild(errorDiv);
    }

    // Register form submission
    registerForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        if (!validateForm()) {
            return;
        }

        registerBtn.classList.add('loading');
        registerBtn.textContent = 'Creating Account...';

        const username = document.getElementById('registerUsername').value;
        const email = document.getElementById('registerEmail').value;
        const password = document.getElementById('registerPassword').value;

        // Prepare registration data
        const registerData = {
            username,
            email,
            password
        };

        // Add exchange keys if provided
        if (enableBinance.checked) {
            registerData.binance_api_key = document.getElementById('binanceApiKey').value;
            registerData.binance_secret_key = document.getElementById('binanceSecretKey').value;
        }

        if (enableSolana.checked) {
            registerData.solana_private_key = document.getElementById('solanaPrivateKey').value;
        }

        try {
            const response = await fetch('/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(registerData),
            });

            const data = await response.json();

            if (response.ok) {
                // Store token and redirect to dashboard
                localStorage.setItem('token', data.token);
                localStorage.setItem('user', JSON.stringify(data.user));

                // Show success message
                showSuccess('Account created successfully! Redirecting to dashboard...');

                // Redirect after a short delay
                setTimeout(() => {
                    window.location.href = '/dashboard';
                }, 500);
            } else {
                showError('registerFormElement', data.error || 'Registration failed');
            }
        } catch (error) {
            showError('registerFormElement', 'Network error: ' + error.message);
        } finally {
            registerBtn.classList.remove('loading');
            registerBtn.textContent = 'Create Account';
        }
    });

    function showSuccess(message) {
        // Remove existing success messages
        document.querySelectorAll('.success-message').forEach(el => el.remove());

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

        registerForm.appendChild(successDiv);
    }
});
