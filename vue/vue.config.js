const { defineConfig } = require('@vue/cli-service');

module.exports = defineConfig({
  transpileDependencies: true,

  // Development server configuration
  devServer: {
    port: 8081, // Run Vue.js on port 8081
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // Proxy to Go server on port 8080
        changeOrigin: true,
        pathRewrite: { '^/api': '' },
      },
    },
  },
});
