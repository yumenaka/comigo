module.exports = {
  "publicPath": "",
  "outputDir": "static",
  "assetsDir": "assets",
  "devServer": {
    "host": "0.0.0.0",
    "port": 48080,
    "disableHostCheck": true,
    "proxy": {
      "/*": {
        "target": "http://localhost:1234",
        "secure": false,
        "changeOrigin": true
      }
    }
  },

  "transpileDependencies": [
    "vuetify"
  ],

  publicPath: '',
  outputDir: 'static',
  assetsDir: 'assets'
}