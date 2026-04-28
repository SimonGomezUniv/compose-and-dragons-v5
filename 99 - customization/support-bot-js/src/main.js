const express = require('express');
const cors = require('cors');
const config = require('./config');
const apiRoutes = require('./api/routes');

const app = express();

// Middleware
app.use(cors());
app.use(express.json());

// Logging
app.use((req, res, next) => {
  console.log(`${new Date().toISOString()} ${req.method} ${req.path}`);
  next();
});

// Routes
app.use('/api', apiRoutes);

// Health endpoint root
app.get('/', (req, res) => {
  res.json({
    service: 'Support Ticket Analysis Bot',
    version: '0.1.0',
    status: 'running',
    endpoints: {
      health: 'GET /api/health',
      analyze: 'POST /api/analyze',
      stats: 'GET /api/stats',
    },
  });
});

// 404 handler
app.use((req, res) => {
  res.status(404).json({ error: 'Endpoint not found' });
});

// Error handler
app.use((err, req, res, next) => {
  console.error('Error:', err);
  res.status(500).json({
    error: 'Internal server error',
    message: err.message,
  });
});

// Démarrer serveur
const PORT = config.api.port;
const HOST = config.api.host;

app.listen(PORT, HOST, () => {
  console.log(`
╔════════════════════════════════════════╗
║   🤖 Support Bot Running               ║
║   Server: http://${HOST}:${PORT}      
║   Node.js: ${process.version}
║   Environment: ${process.env.NODE_ENV || 'development'}
╚════════════════════════════════════════╝
  `);
});

module.exports = app;
