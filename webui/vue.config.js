module.exports = {
  "publicPath": "",
  "outputDir": "static/resources",

  "devServer": {
    "host": "0.0.0.0",
    "port": 8080,
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
  assetsDir: ''
}