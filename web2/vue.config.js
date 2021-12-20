module.exports = {
  // 在浏览器中展示源代码
  configureWebpack: {
    devtool: 'source-map',
  },

  devServer: {
    host: "0.0.0.0",
    port: 4080,
    disableHostCheck: true,
    proxy: {
      '/': {
        //后端服务器地址
        "target": "http://192.168.3.219:1234",
        //允许跨域
        "changeOrigin": true
      }
    }
  },
  publicPath: '',
  outputDir: '../routers/static',
  assetsDir: 'assets'
}