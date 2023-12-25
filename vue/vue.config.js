const { defineConfig } = require('@vue/cli-service');

module.exports = defineConfig({
  transpileDependencies: true,

  // Development server configuration
  devServer: {
    proxy: {
      // Proxy all requests starting with /api to your Echo server
      '/api': {
        target: 'http://localhost:8080', // Change this to your Echo server's URL
        changeOrigin: true,
        pathRewrite: { '^/api': '' }, // Remove '/api' prefix when forwarding to Echo server
      },
    },
  },
});
