const express = require('express');
const app = express();
const port = 80;

app.get('/', (req, res) => {
  res.json({
    message: "Test aplikasi yang di-deploy via CI/CD otomatis!",
    status: "Database connection (simulated) is OK!",
    db_host: process.env.DB_HOST || "Tidak ada koneksi DB"
  });
});

app.listen(port, () => {
  console.log(`App running on port ${port}`);
});
