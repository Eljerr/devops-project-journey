const express = require('express');
const mysql = require('mysql2/promise');
const rateLimit = require('express-rate-limit');
const app = express();

// Rate limiter untuk membatasi request ke endpoint yang melakukan query database (mencegah DoS)
const dbLimiter = rateLimit({
  windowMs: 1 * 60 * 1000, // 1 menit
  max: 60, // Batasi maksimal 60 request per IP per menit
  message: {
    status: 429,
    message: "Too many requests to database health check, please try again later."
  },
  standardHeaders: true,
  legacyHeaders: false,
});

app.get('/', (req, res) => {
  res.json({
    message: "Test aplikasi yang di-deploy via CI/CD otomatis dan! 3x",
    status: "Database connection (simulated) is OK!",
    db_host: process.env.DB_HOST || "Tidak ada koneksi DB"
  });
});

app.get('/healtz', (req,res) => {
  res.status(200).send('OK')
})

app.get('/ready', dbLimiter, async (req, res) => {
  try {
    const connection = await mysql.createConnection({
      host: process.env.DB_HOST || 'mysql-service',
      user: process.env.DB_USER || 'root',
      password: process.env.DB_PASSWORD || 'secret',
      database: process.env.DB_NAME || 'mysql',
      port: process.env.DB_PORT || 3306
    });
    await connection.query('SELECT 1');
    console.log('✅ Database connection successful!');
    await connection.end();
    res.status(200).json({ status: 'ready', db: 'connected' });
  } catch (error) {
    console.error('❌ Database connection failed:', error.message);
    res.status(503).json({ status: 'not ready', db: 'disconnected', error: error.message });
  }
})

module.exports = app
