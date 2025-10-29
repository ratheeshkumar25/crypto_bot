const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const path = require('path');

const app = express();
const PORT = 8080;
const API_BASE_URL = 'http://localhost:3000'; // Backend API URL

// Middleware
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(express.static(path.join(__dirname, 'public')));
app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));

// Proxy middleware
app.use('/api', createProxyMiddleware({ target: API_BASE_URL, changeOrigin: true, logLevel: 'debug' }));

// Routes
app.get('/', (req, res) => {
    res.render('landing', { title: 'CryptoHack Bot - Welcome' });
});

app.get('/register', (req, res) => {
    res.render('register', { title: 'Register - CryptoHack Bot' });
});

app.get('/login', (req, res) => {
    res.render('login', { title: 'Login - CryptoHack Bot' });
});

app.get('/dashboard', (req, res) => {
    res.render('index', { title: 'CryptoHack Bot - Dashboard' });
});

app.get('/test', (req, res) => {
    res.json({ message: 'test successful' });
});


// Start server
app.listen(PORT, () => {
    console.log(`Frontend server running at http://localhost:${PORT}`);
});
