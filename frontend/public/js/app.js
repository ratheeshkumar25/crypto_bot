// Frontend JavaScript for Crypto Trading Bot

document.addEventListener('DOMContentLoaded', function() {
    // Authentication elements
    const showLoginBtn = document.getElementById('showLoginBtn');
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const loginFormElement = document.getElementById('loginFormElement');
    const registerFormElement = document.getElementById('registerFormElement');
    const userInfo = document.getElementById('userInfo');
    const usernameDisplay = document.getElementById('usernameDisplay');
    const logoutBtn = document.getElementById('logoutBtn');
    const showRegister = document.getElementById('showRegister');
    const showLogin = document.getElementById('showLogin');

    // Trading elements
    const priceForm = document.getElementById('priceForm');
    const predictForm = document.getElementById('predictForm');
    const priceResult = document.getElementById('priceResult');
    const predictResult = document.getElementById('predictResult');
    const startChartBtn = document.getElementById('startChart');
    const stopChartBtn = document.getElementById('stopChart');
    const chartSymbolSelect = document.getElementById('chartSymbol');
    const chartStrategySelect = document.getElementById('chartStrategy');
    const currentPriceDiv = document.getElementById('currentPriceSplitView');
    const showGuideBtn = document.getElementById('showGuide');
    const closeGuideBtn = document.getElementById('closeGuide');
    const strategyGuide = document.getElementById('strategyGuide');

    let token = localStorage.getItem('token');
    let currentUser = JSON.parse(localStorage.getItem('user')) || null;

    console.log('user from local storage:', localStorage.getItem('user'));
    console.log('Token:', token);
    console.log('Current User:', currentUser);

    // Initialize authentication state
    updateAuthUI();







    let chart;
    let chartInterval;
    let priceData = [];
    let timeLabels = [];
    let signalDatasets = [];
    let wsConnection = null;

    // Initialize Chart.js
    const ctx = document.getElementById('priceChart').getContext('2d');
    chart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: timeLabels,
            datasets: [{
                label: 'Price',
                data: priceData,
                borderColor: 'rgb(75, 192, 192)',
                tension: 0.1,
                pointRadius: 0,
                borderWidth: 2
            }]
        },
        options: {
            responsive: true,
            scales: {
                x: {
                    display: true,
                    title: {
                        display: true,
                        text: 'Time'
                    }
                },
                y: {
                    display: true,
                    title: {
                        display: true,
                        text: 'Price (USDT)'
                    }
                }
            },
            plugins: {
                legend: {
                    display: true
                }
            }
        }
    });

    // Start live chart
    startChartBtn.addEventListener('click', function() {
        if (chartInterval) clearInterval(chartInterval);
        if (wsConnection) wsConnection.close();

        const symbol = chartSymbolSelect.value;
        const strategy = chartStrategySelect.value;

        // Start WebSocket connection for real-time updates
        startWebSocket(symbol, strategy);
    });

    // Stop live chart
    stopChartBtn.addEventListener('click', function() {
        if (chartInterval) {
            clearInterval(chartInterval);
            chartInterval = null;
        }
        if (wsConnection) {
            wsConnection.close();
            wsConnection = null;
        }
        // Clear signal datasets
        clearSignalDatasets();
    });

    // Function to clear signal datasets
    function clearSignalDatasets() {
        // Remove all datasets except the first one (price data)
        while (chart.data.datasets.length > 1) {
            chart.data.datasets.pop();
        }
        signalDatasets = [];
        chart.update();
    }

    // Function to add signal datasets
    // Function to add signal datasets
    async function addSignalDatasets(symbol, strategy) {
        if (!strategy) return;

        try {
            const response = await fetch(`/api/signals/${strategy}?symbol=${symbol}`, {
                headers: getAuthHeaders()
            });
            const data = await response.json();

            if (response.ok && data.signals) {
                clearSignalDatasets();

                data.signals.forEach(signal => {
                    const isBuy = signal.Type === 'BUY';
                    const mainColor = isBuy ? 'rgba(34, 197, 94, 0.8)' : 'rgba(239, 68, 68, 0.8)';
                    const profitColor = 'rgba(74, 222, 128, 0.2)'; // Light green for profit
                    const lossColor = 'rgba(252, 165, 165, 0.2)';   // Light red for loss

                    // Main signal line (Buy or Sell)
                    const mainDataset = {
                        label: `${signal.Type} at $${signal.Price.toFixed(2)}`,
                        data: Array(timeLabels.length).fill(signal.Price),
                        borderColor: mainColor,
                        borderWidth: 2,
                        pointRadius: 0,
                        fill: false,
                    };
                    chart.data.datasets.push(mainDataset);

                    // Take Profit line
                    const takeProfitDataset = {
                        label: `Take Profit at $${signal.TakeProfit.toFixed(2)}`,
                        data: Array(timeLabels.length).fill(signal.TakeProfit),
                        borderColor: 'rgba(34, 197, 94, 0.7)',
                        borderDash: [5, 5],
                        borderWidth: 1.5,
                        pointRadius: 0,
                        fill: false,
                    };
                    chart.data.datasets.push(takeProfitDataset);

                    // Stop Loss line
                    const stopLossDataset = {
                        label: `Stop Loss at $${signal.StopLoss.toFixed(2)}`,
                        data: Array(timeLabels.length).fill(signal.StopLoss),
                        borderColor: 'rgba(239, 68, 68, 0.7)',
                        borderDash: [5, 5],
                        borderWidth: 1.5,
                        pointRadius: 0,
                        fill: false,
                    };
                    chart.data.datasets.push(stopLossDataset);

                    // Shaded profit/loss zones
                    const profitZoneDataset = {
                        label: 'Profit Zone',
                        data: Array(timeLabels.length).fill(signal.TakeProfit),
                        borderColor: 'transparent',
                        backgroundColor: isBuy ? profitColor : lossColor,
                        pointRadius: 0,
                        fill: {
                            target: { value: signal.Price },
                            above: isBuy ? profitColor : 'transparent',
                            below: !isBuy ? profitColor : 'transparent'
                        },
                    };
                    chart.data.datasets.push(profitZoneDataset);

                    const lossZoneDataset = {
                        label: 'Loss Zone',
                        data: Array(timeLabels.length).fill(signal.StopLoss),
                        borderColor: 'transparent',
                        backgroundColor: isBuy ? lossColor : profitColor,
                        pointRadius: 0,
                        fill: {
                            target: { value: signal.Price },
                            above: !isBuy ? lossColor : 'transparent',
                            below: isBuy ? lossColor : 'transparent'
                        },
                    };
                    chart.data.datasets.push(lossZoneDataset);
                });

                chart.update();
            }
        } catch (error) {
            console.error('Error fetching signals:', error);
        }
    }

    // Function to start WebSocket connection
    function startWebSocket(symbol, strategy) {
        const wsUrl = `ws://localhost:8080/ws/price?symbol=${symbol}`;
        wsConnection = new WebSocket(wsUrl);

        wsConnection.onopen = function(event) {
            console.log('WebSocket connection opened');
            // Add signals if strategy is selected
            if (strategy) {
                addSignalDatasets(symbol, strategy);
            }
        };

        wsConnection.onmessage = function(event) {
            try {
                const data = JSON.parse(event.data);

                if (data.error) {
                    throw new Error(data.error);
                }

                const now = new Date().toLocaleTimeString();
                timeLabels.push(now);
                priceData.push(data.price);

                // Keep only last 50 data points
                if (timeLabels.length > 50) {
                    timeLabels.shift();
                    priceData.shift();
                }

                chart.update();
                const bidPrice = data.price * 0.9995;
                const askPrice = data.price * 1.0005;
                currentPriceDiv.innerHTML = `
                    <div class="price-panel bid">
                        <h3>Bid Price</h3>
                        <div class="price">$${bidPrice.toFixed(2)}</div>
                    </div>
                    <div class="price-panel ask">
                        <h3>Ask Price</h3>
                        <div class="price">$${askPrice.toFixed(2)}</div>
                    </div>
                `;
            } catch (error) {
                currentPriceDiv.innerHTML = `
                    <div class="result error">
                        <strong>Error updating chart:</strong> ${error.message}
                    </div>
                `;
            }
        };

        wsConnection.onclose = function(event) {
            console.log('WebSocket connection closed');
        };

        wsConnection.onerror = function(error) {
            console.error('WebSocket error:', error);
            currentPriceDiv.innerHTML = `
                <div class="result error">
                    <strong>WebSocket connection failed. Falling back to polling.</strong>
                </div>
            `;
            // Fallback to polling
            updateChartPolling(symbol, strategy);
        };
    }

    // Function to update chart with polling (fallback)
    async function updateChartPolling(symbol, strategy) {
        try {
            const response = await fetch(`/price?exchange=binance&symbol=${symbol}`);
            const data = await response.json();

            if (response.ok) {
                const now = new Date().toLocaleTimeString();
                timeLabels.push(now);
                priceData.push(data.price);

                // Keep only last 50 data points
                if (timeLabels.length > 50) {
                    timeLabels.shift();
                    priceData.shift();
                }

                // Add signals if strategy is selected
                if (strategy) {
                    await addSignalDatasets(symbol, strategy);
                }

                chart.update();
                currentPriceDiv.innerHTML = `
                    <div class="result success">
                        <strong>Current Price for ${symbol}:</strong> $${data.price}<br>
                        <small>Last updated: ${now}</small>
                    </div>
                `;
            } else {
                throw new Error(data.error || 'Failed to get price');
            }
        } catch (error) {
            currentPriceDiv.innerHTML = `
                <div class="result error">
                    <strong>Error updating chart:</strong> ${error.message}
                </div>
            `;
        }
    }

    // Handle price form submission
    priceForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const exchange = document.getElementById('exchange').value;
        const symbol = document.getElementById('symbol').value;

        try {
            const response = await fetch(`/price?exchange=${exchange}&symbol=${symbol}`, {
                headers: getAuthHeaders()
            });
            const data = await response.json();

            if (response.ok) {
                priceResult.innerHTML = `
                    <div class="result success">
                        <strong>Current Price for ${symbol}:</strong><br>
                        Price: $${data.price}<br>
                        Timestamp: ${new Date(data.timestamp).toLocaleString()}
                    </div>
                `;
            } else {
                throw new Error(data.error || 'Failed to get price');
            }
        } catch (error) {
            priceResult.innerHTML = `
                <div class="result error">
                    <strong>Error:</strong> ${error.message}
                </div>
            `;
        }
    });

    // Handle prediction form submission
    predictForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const strategy = document.getElementById('strategy').value;
        const symbol = document.getElementById('predictSymbol').value;
        const investment = document.getElementById('investment').value;
        const timeframe = document.getElementById('timeframe').value;

        try {
            const response = await fetch(`/predict?strategy=${strategy}&symbol=${symbol}&investment=${investment}&timeframe=${timeframe}`, {
                headers: getAuthHeaders()
            });
            const data = await response.json();

            if (response.ok) {
                predictResult.innerHTML = `
                    <div class="result success">
                        <strong>Profit Prediction for ${strategy.toUpperCase()} Strategy:</strong><br>
                        Symbol: ${symbol}<br>
                        Investment: $${investment}<br>
                        Timeframe: ${timeframe}<br>
                        Predicted Profit: $${data.predictedProfit.toFixed(2)}<br>
                        Profit Percentage: ${data.profitPercentage.toFixed(2)}%
                    </div>
                `;
            } else {
                throw new Error(data.error || 'Failed to get prediction');
            }
        } catch (error) {
            predictResult.innerHTML = `
                <div class="result error">
                    <strong>Error:</strong> ${error.message}
                </div>
            `;
        }
    });

    // Show strategy guide
    showGuideBtn.addEventListener('click', function() {
        strategyGuide.style.display = 'flex';
    });

    // Close strategy guide
    closeGuideBtn.addEventListener('click', function() {
        strategyGuide.style.display = 'none';
    });

    // Authentication event listeners
    showLoginBtn.addEventListener('click', function() {
        loginForm.style.display = 'block';
        registerForm.style.display = 'none';
        showLoginBtn.style.display = 'none';
    });

    showRegister.addEventListener('click', function() {
        loginForm.style.display = 'none';
        registerForm.style.display = 'block';
    });

    showLogin.addEventListener('click', function() {
        registerForm.style.display = 'none';
        loginForm.style.display = 'block';
    });

    logoutBtn.addEventListener('click', function() {
        logout();
    });

    // Login form submission
    loginFormElement.addEventListener('submit', async function(e) {
        e.preventDefault();
        const username = document.getElementById('loginUsername').value;
        const password = document.getElementById('loginPassword').value;

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
                token = data.token;
                currentUser = data.user;
                localStorage.setItem('token', token);
                localStorage.setItem('user', JSON.stringify(data.user));
                updateAuthUI();
                window.location.href = '/dashboard';
            } else {
                alert(data.error || 'Login failed');
            }
        } catch (error) {
            alert('Login error: ' + error.message);
        }
    });

    // Register form submission
    registerFormElement.addEventListener('submit', async function(e) {
        e.preventDefault();
        const username = document.getElementById('registerUsername').value;
        const email = document.getElementById('registerEmail').value;
        const password = document.getElementById('registerPassword').value;

        try {
                                    const response = await fetch('/api/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, email, password }),
            });

            const data = await response.json();

            if (response.ok) {
                token = data.token;
                currentUser = data.user;
                localStorage.setItem('token', token);
                localStorage.setItem('user', JSON.stringify(data.user));
                updateAuthUI();
                registerForm.style.display = 'none';
                setTimeout(() => {
                    window.location.href = '/dashboard';
                }, 100);
            } else {
                alert(data.error || 'Registration failed');
            }
        } catch (error) {
            alert('Registration error: ' + error.message);
        }
    });

    // Authentication helper functions
    function updateAuthUI() {
        console.log('updateAuthUI called - token:', token, 'currentUser:', currentUser);
        if (token && currentUser) {
            console.log('Showing user info');
            if (userInfo) userInfo.style.display = 'block';
            if (usernameDisplay) usernameDisplay.textContent = `Welcome, ${currentUser.username}!`;
            if (logoutBtn) logoutBtn.style.display = 'inline-block';
            if (showLoginBtn) showLoginBtn.style.display = 'none';
            if (loginForm) loginForm.style.display = 'none';
            if (registerForm) registerForm.style.display = 'none';
            addUserTradesToChart();
        } else {
            console.log('Showing login button');
            if (userInfo) userInfo.style.display = 'none';
            if (logoutBtn) logoutBtn.style.display = 'none';
            if (showLoginBtn) showLoginBtn.style.display = 'block';
            if (loginForm) loginForm.style.display = 'none';
            if (registerForm) registerForm.style.display = 'none';
        }
    }

    function logout() {
        token = null;
        currentUser = null;
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        updateAuthUI();
        setTimeout(() => {
            window.location.href = '/';
        }, 100);
    }



    // Helper function to get auth headers
    function getAuthHeaders() {
        return token ? { 'Authorization': `Bearer ${token}` } : {};
    }

    // Function to add user trades to the chart
    async function addUserTradesToChart() {
        if (!token) return;

        try {
            const response = await fetch('/api/trades', {
                headers: getAuthHeaders()
            });
            const data = await response.json();

            if (response.ok && data.trades) {
                data.trades.forEach(trade => {
                    const isBuy = trade.Side === 'BUY';
                    const color = isBuy ? 'blue' : 'purple';

                    // Entry price
                    const entryDataset = {
                        label: `${trade.Side} @ ${trade.Price.toFixed(2)}`,
                        data: Array(timeLabels.length).fill(trade.Price),
                        borderColor: color,
                        borderWidth: 2,
                        pointRadius: 0,
                        fill: false,
                    };
                    chart.data.datasets.push(entryDataset);

                    // Take profit
                    if (trade.TakeProfit > 0) {
                        const takeProfitDataset = {
                            label: `TP @ ${trade.TakeProfit.toFixed(2)}`,
                            data: Array(timeLabels.length).fill(trade.TakeProfit),
                            borderColor: color,
                            borderDash: [5, 5],
                            borderWidth: 1.5,
                            pointRadius: 0,
                            fill: false,
                        };
                        chart.data.datasets.push(takeProfitDataset);
                    }

                    // Stop loss
                    if (trade.StopLoss > 0) {
                        const stopLossDataset = {
                            label: `SL @ ${trade.StopLoss.toFixed(2)}`,
                            data: Array(timeLabels.length).fill(trade.StopLoss),
                            borderColor: color,
                            borderDash: [10, 10],
                            borderWidth: 1.5,
                            pointRadius: 0,
                            fill: false,
                        };
                        chart.data.datasets.push(stopLossDataset);
                    }
                });
                chart.update();
            }
        } catch (error) {
            console.error('Error fetching user trades:', error);
        }
    }
});
