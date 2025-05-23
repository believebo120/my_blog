const { defineConfig } = require('@vue/cli-service');
const NodePolyfillPlugin = require('node-polyfill-webpack-plugin');

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8081,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        pathRewrite: { '^/api': '' }
      },
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  chainWebpack: (config) => {
    // 添加 Node.js 模块的 polyfill
    config.plugin('node-polyfill').use(NodePolyfillPlugin);
    // 配置 fallback（修正此处）
    config.resolve.set('fallback', {
      os: require.resolve('os-browserify/browser')
    });
  }
});