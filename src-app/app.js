const express = require('express');
const app = express();

app.get('/', (req, res) => {
  res.json({
    message: "Test aplikasi yang di-deploy via CI/CD otomatis dan! 1x",
    status: "Database connection (simulated) is OK!",
    db_host: process.env.DB_HOST || "Tidak ada koneksi DB"
  });
});
module.exports = app
