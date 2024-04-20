const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/submit', // Ubah '/api' sesuai dengan path endpoint backend Anda
    createProxyMiddleware({
      target: 'http://localhost:8080', // Ganti dengan URL server backend Go Anda
      changeOrigin: true,
    })
  );
};
