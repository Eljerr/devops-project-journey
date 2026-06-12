const express = require('express');
const app = express();

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

app.get('/ready', async (req,res) => {
  try {
    const connection = await mysql.createConnection({
      host: 'localhost',
      user: 'root',
      password: 'secret',
      database: 'mysql'
    })
      await connection.query('SELECT 1')
      console.log('✅ Database connection successful!')
      
      await connection.end()
  } catch (error) {
    console.error('Database connection failed')
  }
})

module.exports = app
