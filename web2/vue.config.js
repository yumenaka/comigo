module.exports = {
  "devServer": {
    "host": "0.0.0.0",
    "port": 48080,
    "disableHostCheck": true,
    "proxy": {
      "/*": {
        "target": "http://192.168.3.219:1234",
        "secure": false,
        "changeOrigin": true
      }
    }
  },
  publicPath: '',
  outputDir: 'static',
  assetsDir: 'assets'
}